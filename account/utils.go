package account

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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

var ApiClient *http.Client

func init() {
	ApiClient = http.DefaultClient
}

func (index Verb) String() string {
	return [...]string{"POST", "GET", "DELETE"}[index]
}

func NewRequestWithHeaders(verb Verb, id uuid.UUID) *http.Request {
	req, err := http.NewRequest(verb.String(), endpointString(id), nil)
	Check(err)

	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Date", time.Now().String())
	req.Header.Add("Accept", "application/vnd.api+json")
	return req
}

func endpointString(id uuid.UUID) string {
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

func HandleResponse(response *http.Response, verb Verb) (*AccountApiResponse, error) {
	var responseWrapper AccountApiResponse
	responseWrapper.Status = &response.Status
	responseWrapper.StatusCode = &response.StatusCode

	responseBody, err := ioutil.ReadAll(response.Body)
	Check(err)
	if response.StatusCode < 400 {
		switch verb {
		case CREATE:
			var createdAccount Account
			Check(json.Unmarshal(responseBody, &createdAccount))
			responseWrapper.ResponseBody = &createdAccount
			return &responseWrapper, nil
		case FETCH:
			var fetchedAccount Account
			Check(json.Unmarshal(responseBody, &fetchedAccount))
			responseWrapper.ResponseBody = &fetchedAccount
			return &responseWrapper, nil
		case DELETE:
			return nil, errors.New("NOT IMPLEMENTED YET")
		default:
			return nil, errors.New("UNHANDLED HTTP VERB")
		}
	} else {
		return nil, errors.New(string(responseBody))
	}
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
