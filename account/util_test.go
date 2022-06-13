package account

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHttpVerbString(t *testing.T) {
	subtests := []struct {
		name      string
		verb      httpMethod
		expString string
	}{
		{
			"Verb POST mapped successfully",
			createMethod,
			"POST",
		},
		{
			"Verb FETCH mapped successfully",
			fetchMethod,
			"GET",
		},
		{
			"Verb DELETE mapped successfully",
			deleteMethod,
			"DELETE",
		},
	}
	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			if subtest.verb.String() != subtest.expString {
				t.Errorf("expected httpVerb %d to produce %s", subtest.verb, subtest.verb.String())
			}
		})
	}
}

func TestNewRequestWithHeaders(t *testing.T) {

	subtests := []struct {
		name            string
		httpNewRequest  func(method, url string, body io.Reader) (*http.Request, error)
		httpVerb        httpMethod
		id              uuid.UUID
		version         *int64
		expHostHdr      string
		expAcceptHdr    string
		expVersionParam string
		expPanic        bool
	}{
		{
			name: "New POST request with Headers",
			httpNewRequest: func(method, httpUrl string, body io.Reader) (*http.Request, error) {
				return &http.Request{
					Header: make(http.Header),
					Method: http.MethodPost,
					URL: &url.URL{
						RawPath: "http://localhost:8080/v1/organization/accounts",
					},
				}, nil
			},
			httpVerb:     createMethod,
			id:           uuid.New(),
			expHostHdr:   Host,
			expAcceptHdr: "application/vnd.api+json",
		},
		{
			name: "New DELETE request with Headers",
			httpNewRequest: func(method, httpUrl string, body io.Reader) (*http.Request, error) {
				return &http.Request{
					Header:     make(http.Header),
					RequestURI: "http://localhost:8080/v1/organization/accounts",
					URL: &url.URL{
						RawPath: "http://localhost:8080/v1/organization/accounts",
					},
				}, nil
			},
			httpVerb:        deleteMethod,
			id:              uuid.New(),
			version:         int64ToPointer(0),
			expVersionParam: "0",
			expHostHdr:      Host,
			expAcceptHdr:    "application/vnd.api+json",
		},
		{
			name: "http.NewRequest returns error and func panics",
			httpNewRequest: func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("Error")
			},
			expPanic: true,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			httpNewRequest = subtest.httpNewRequest
			if subtest.expPanic {
				assertPanicNewReq(t, newRequestWithHeaders)
			} else {
				result := newRequestWithHeaders(subtest.httpVerb, subtest.id, subtest.version)
				if result.Header.Get("Host") != subtest.expHostHdr {
					t.Errorf("expected header Host (%+v), got (%+v)", subtest.expHostHdr, result.Header.Get("Host"))
				}
				if result.Header.Get("Accept") != subtest.expAcceptHdr {
					t.Errorf("expected header Accept (%+v), got (%+v)", subtest.expAcceptHdr, result.Header.Get("Accept"))
				}
				testDateHdr(result.Header.Get("Date"), t)
				versionQueryParam := result.URL.Query().Get("version")
				if versionQueryParam != subtest.expVersionParam {
					t.Errorf("expected version param (%+v), got (%+v)", subtest.expVersionParam, versionQueryParam)
				}
			}
		})
	}
}

func testDateHdr(dateHdrStr string, t *testing.T) {
	dateHdrTime, err := time.Parse(time.RFC3339Nano, dateHdrStr)
	if err != nil {
		t.Errorf("got error from Date header conversion %s", err.Error())
	}
	timeNow := time.Now()
	diff := timeNow.Sub(dateHdrTime)
	if diff*time.Second < 1 {
		t.Errorf("expected header Date within 1 sec before (%s), got (%s)", timeNow.Format(time.RFC3339Nano), dateHdrStr)
	}
}

func int64ToPointer(i int64) *int64 {
	return &i
}

func assertPanicNewReq(t *testing.T, f func(verb httpMethod, id uuid.UUID, version *int64) *http.Request) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f(createMethod, uuid.Nil, nil)
}
func TestHandleResponse(t *testing.T) {
	subtests := []struct {
		name                string
		httpVerb            httpMethod
		responseParam       *http.Response
		handleCreateOrFetch func(responseBody []byte, responseWrapper AccountApiResponse, verb httpMethod) (*AccountApiResponse, error)
		handleDelete        func(responseWrapper AccountApiResponse, responseBody []byte) (*AccountApiResponse, error)
		readRespBody        func(r io.Reader) ([]byte, error)
		expResponse         *AccountApiResponse
		expError            error
		expPanic            bool
	}{
		{
			name:     "Handle successful POST response",
			httpVerb: createMethod,
			responseParam: &http.Response{
				Status:     "Created",
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewBufferString("Successful POST resp")),
			},
			handleCreateOrFetch: func(responseBody []byte, responseWrapper AccountApiResponse, verb httpMethod) (*AccountApiResponse, error) {
				return exp_res_created_success, nil
			},
			expResponse: exp_res_created_success,
			readRespBody: func(r io.Reader) ([]byte, error) {
				return []byte("Successful POST resp"), nil
			},
		},
		{
			name:     "Handle successful FETCH response",
			httpVerb: fetchMethod,
			responseParam: &http.Response{
				Status:     "OK",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("Successful Fetch Res")),
			},
			handleCreateOrFetch: func(responseBody []byte, responseWrapper AccountApiResponse, verb httpMethod) (*AccountApiResponse, error) {
				return exp_res_fetch_success, nil
			},
			expResponse: exp_res_fetch_success,
			readRespBody: func(r io.Reader) ([]byte, error) {
				return []byte("Successful Fetch resp"), nil
			},
		},
		{
			name:     "Handle successful Delete response",
			httpVerb: deleteMethod,
			responseParam: &http.Response{
				Status:     "No content",
				StatusCode: http.StatusNoContent,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			handleDelete: func(responseWrapper AccountApiResponse, responseBody []byte) (*AccountApiResponse, error) {
				return exp_res_deleted_success, nil
			},
			expResponse: exp_res_deleted_success,
			readRespBody: func(r io.Reader) ([]byte, error) {
				return []byte(""), nil
			},
		},
		{
			name:     "Handle Status BAD REQUEST 400",
			httpVerb: createMethod,
			responseParam: &http.Response{
				Status:     "Bad request",
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			expResponse: nil,
			readRespBody: func(r io.Reader) ([]byte, error) {
				return []byte(""), nil
			},
			expError: &ApiError{
				StatusCode:   http.StatusBadRequest,
				Status:       "Bad request",
				ResponseBody: "",
				Message:      "GOT ERROR STATUS CODE OF 400, STATUS Bad request",
			},
		},
		{
			name:     "ioutil.ReadAll (readRespBody) returns error and func panics",
			expPanic: true,
			readRespBody: func(r io.Reader) ([]byte, error) {
				return nil, errors.New("Unable to read Response Body")
			},
			responseParam: &http.Response{
				Status:     "Created",
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewBufferString("Successful POST resp")),
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			readRespBody = subtest.readRespBody
			handleCreateOrFetch = subtest.handleCreateOrFetch
			handleDelete = subtest.handleDelete
			if subtest.expPanic {
				assertPanicHandleResp(t, handleResponse, subtest.responseParam)
			} else {
				result, err := handleResponse(subtest.responseParam, subtest.httpVerb)
				if !errors.Is(err, subtest.expError) {
					t.Errorf("expected error (%v), got error (%v)", subtest.expError, err)
				}
				if !reflect.DeepEqual(result, subtest.expResponse) {
					t.Errorf("expected (%+v), got (%+v)", subtest.expResponse, result)
				}
			}
		})
	}
}

func assertPanicHandleResp(t *testing.T, f func(response *http.Response, verb httpMethod) (*AccountApiResponse, error), resp *http.Response) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f(resp, createMethod)
}
