package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const Url string = "http://localhost:8080/v1/organisation/accounts"

type Verb int

const (
	CREATE Verb = iota
	FETCH
	DELETE
)

const (
	API_ERROR_FORMATTING         = "ACCOUNT API ERROR\nSTATUS CODE : %d\nSTATUS : %s\nRESPONSE BODY : %s\nMESSAGE : %s"
	DELETE_INCORRECT_STATUS_CODE = "DELETE OPERATION GOT INCORRECT STATUS CODE. EXPECTED: %d, GOT: %d"
	ERROR_STATUS_CODE            = "GOT ERROR STATUS CODE OF %d\n, STATUS %s"
)

var ApiClient *http.Client

func init() {
	ApiClient = http.DefaultClient
}

func (index Verb) String() string {
	return [...]string{"POST", "GET", "DELETE"}[index]
}

func newRequestWithHeaders(verb Verb, id uuid.UUID, version *int64) *http.Request {
	req, err := http.NewRequest(verb.String(), endpointString(verb, id), nil)
	check(err)
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Date", time.Now().String())
	req.Header.Add("Accept", "application/vnd.api+json")

	if verb == DELETE {
		q := req.URL.Query()
		q.Add("version", strconv.Itoa(int(*version)))
		req.URL.RawQuery = q.Encode()
	}
	return req
}

func endpointString(verb Verb, id uuid.UUID) string {
	if id == uuid.Nil {
		return Url
	} else {
		var sb strings.Builder
		sb.WriteString(Url)
		sb.WriteString("/")
		sb.WriteString(id.String())
		return sb.String()
	}
}

func handleResponse(response *http.Response, verb Verb) (*AccountApiResponse, error) {
	var responseWrapper AccountApiResponse
	responseWrapper.Status = &response.Status
	responseWrapper.StatusCode = &response.StatusCode

	responseBody, err := ioutil.ReadAll(response.Body)
	check(err)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		switch verb {
		case CREATE, FETCH:
			return handleCreateOrFetchResponse(responseBody, responseWrapper, verb)
		case DELETE:
			return handleDeleteResponse(responseWrapper, responseBody)
		default:
			return nil, errors.New("UNHANDLED HTTP VERB")
		}
	} else {
		return nil, &ApiError{*responseWrapper.StatusCode, *responseWrapper.Status, string(responseBody),
			fmt.Sprintf(ERROR_STATUS_CODE, *responseWrapper.StatusCode, *responseWrapper.Status)}
	}
}

func handleCreateOrFetchResponse(responseBody []byte, responseWrapper AccountApiResponse, verb Verb) (*AccountApiResponse, error) {
	switch verb {
	case CREATE:
		if *responseWrapper.StatusCode != http.StatusCreated {
			return nil, &ApiError{*responseWrapper.StatusCode, *responseWrapper.Status, string(responseBody),
				fmt.Sprintf(DELETE_INCORRECT_STATUS_CODE, http.StatusNoContent, *responseWrapper.StatusCode)}
		}
	case FETCH:
		if *responseWrapper.StatusCode != http.StatusOK {
			return nil, &ApiError{*responseWrapper.StatusCode, *responseWrapper.Status, string(responseBody),
				fmt.Sprintf(DELETE_INCORRECT_STATUS_CODE, http.StatusNoContent, *responseWrapper.StatusCode)}
		}
	default:
		return nil, errors.New("HANDLE CREATE OR FETCH FUNCTION CALLED WITH INCORRECT HTTP VERB")
	}
	var account Account
	check(json.Unmarshal(responseBody, &account))
	responseWrapper.ResponseBody = &account
	return &responseWrapper, nil
}

func handleDeleteResponse(responseWrapper AccountApiResponse, responseBody []byte) (*AccountApiResponse, error) {
	if *responseWrapper.StatusCode == http.StatusNoContent {
		return &responseWrapper, nil
	} else {
		return nil, &ApiError{*responseWrapper.StatusCode, *responseWrapper.Status, string(responseBody),
			fmt.Sprintf(DELETE_INCORRECT_STATUS_CODE, http.StatusNoContent, *responseWrapper.StatusCode)}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type ApiError struct {
	StatusCode   int
	Status       string
	ResponseBody string
	Message      string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf(API_ERROR_FORMATTING,
		e.StatusCode, e.Status, e.ResponseBody, e.Message)
}

// func fatalApiError(statusCode int, status string, body string, message string) {
// 	l := log.New(os.Stderr, "", 1)
// 	l.Println(message)
// 	l.Fatalf("Status Code : %d, Status : %s", statusCode, status)
// 	l.Printf("Response Body : %#s\n", body)
// 	os.Exit(1)
// }
