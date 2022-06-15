package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var (
	// Host must contain API http scheme + actual hostname including possible port.
	// It should end with forward slash.
	// Example: "http://localhost:8080/"
	Host string = "http://localhost:8080/"
	// ApiVersion contains api version.
	// It will be concatenated after the Host variable, therefore
	// it shouldn't have a slash prefix, but should have a slash postfix.
	// Example: "v1/"
	ApiVersion string = "v1/"
	// ApiClient is a pointer to the standard http client which will execute the http requests.
	// Has a timeout of 1 second.
	ApiClient *http.Client = &http.Client{Timeout: time.Duration(1) * time.Second}
)

var (
	apiCall             = ApiClient.Do
	jsonMarshal         = json.Marshal
	jsonUnmarshal       = json.Unmarshal
	httpNewRequest      = http.NewRequest
	readRespBody        = ioutil.ReadAll
	newReq              = newRequestWithHeaders
	handleRes           = handleResponse
	handleCreateOrFetch = handleResponseForCreateOrFetch
	handleDelete        = handleDeleteResponse
)

type httpMethod int

const (
	createMethod httpMethod = iota
	fetchMethod
	deleteMethod
	accountsEndpoint                          = "organisation/accounts"
	api_error_formatting                      = "ACCOUNT API ERROR\nSTATUS CODE : %d\nSTATUS : %s\nRESPONSE BODY : %s\nMESSAGE : %s"
	delete_incorrect_status_code_formatting   = "DELETE OPERATION GOT INCORRECT STATUS CODE. EXPECTED: %d, GOT: %d"
	create_incorrect_status_code_formatting   = "CREATE OPERATION GOT INCORRECT STATUS CODE. EXPECTED: %d, GOT: %d"
	fetch_incorrect_status_code_formatting    = "FETCH OPERATION GOT INCORRECT STATUS CODE. EXPECTED: %d, GOT: %d"
	error_status_code_formatting              = "GOT ERROR STATUS CODE OF %d, STATUS %s"
	create_or_fetch_incorrect_verb_formatting = "HANDLE CREATE OR FETCH FUNCTION CALLED WITH INCORRECT HTTP VERB"
)

func (index httpMethod) String() string {
	return [...]string{"POST", "GET", "DELETE"}[index]
}

var newRequestWithHeaders = func(verb httpMethod, id uuid.UUID, version *int64) (*http.Request, error) {
	req, err := httpNewRequest(verb.String(), endpointString(id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Host", Host)
	req.Header.Add("Date", time.Now().Format(time.RFC3339Nano))
	req.Header.Add("Accept", "application/vnd.api+json")

	if verb == deleteMethod {
		q := req.URL.Query()
		q.Add("version", strconv.Itoa(int(*version)))
		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}

var endpointString = func(id uuid.UUID) string {
	finalEndpoint := Host + ApiVersion + accountsEndpoint
	if id == uuid.Nil {
		return finalEndpoint
	} else {
		return finalEndpoint + "/" + id.String()
	}
}

var handleResponse = func(response *http.Response, verb httpMethod) (*AccountApiResponse, error) {
	var responseWrapper AccountApiResponse
	responseWrapper.Status = response.Status
	responseWrapper.StatusCode = response.StatusCode

	responseBody, err := readRespBody(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < http.StatusBadRequest {
		switch verb {
		case createMethod, fetchMethod:
			return handleCreateOrFetch(responseBody, responseWrapper, verb)
		case deleteMethod:
			return handleDelete(responseWrapper, responseBody)
		default:
			return nil, errors.New("UNHANDLED HTTP VERB")
		}
	} else {
		return nil, &ApiError{responseWrapper.StatusCode, responseWrapper.Status, string(responseBody),
			fmt.Sprintf(error_status_code_formatting, responseWrapper.StatusCode, responseWrapper.Status)}
	}
}

var handleResponseForCreateOrFetch = func(responseBody []byte, responseWrapper AccountApiResponse, verb httpMethod) (*AccountApiResponse, error) {
	switch verb {
	case createMethod:
		if responseWrapper.StatusCode != http.StatusCreated {
			return nil, &ApiError{responseWrapper.StatusCode, responseWrapper.Status, string(responseBody),
				fmt.Sprintf(create_incorrect_status_code_formatting, http.StatusCreated, responseWrapper.StatusCode)}
		}
	case fetchMethod:
		if responseWrapper.StatusCode != http.StatusOK {
			return nil, &ApiError{responseWrapper.StatusCode, responseWrapper.Status, string(responseBody),
				fmt.Sprintf(fetch_incorrect_status_code_formatting, http.StatusOK, responseWrapper.StatusCode)}
		}
	default:
		return nil, errors.New(create_or_fetch_incorrect_verb_formatting)
	}
	var account Account
	err := jsonUnmarshal(responseBody, &account)
	if err != nil {
		return nil, err
	}
	responseWrapper.ResponseBody = &account
	return &responseWrapper, nil
}

var handleDeleteResponse = func(responseWrapper AccountApiResponse, responseBody []byte) (*AccountApiResponse, error) {
	if responseWrapper.StatusCode == http.StatusNoContent {
		return &responseWrapper, nil
	} else {
		return nil, &ApiError{responseWrapper.StatusCode, responseWrapper.Status, string(responseBody),
			fmt.Sprintf(delete_incorrect_status_code_formatting, http.StatusNoContent, responseWrapper.StatusCode)}
	}
}

// ApiError is a custom error being returned in case of an error response from form3 API.
type ApiError struct {
	StatusCode   int
	Status       string
	ResponseBody string
	Message      string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf(api_error_formatting,
		e.StatusCode, e.Status, e.ResponseBody, e.Message)
}

func (e *ApiError) Is(tgt error) bool {
	return e.Error() == tgt.Error()
}
