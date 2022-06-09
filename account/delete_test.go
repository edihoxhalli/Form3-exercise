package account

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

var (
	exp_res_deleted_success = &AccountApiResponse{
		ResponseBody: nil,
		StatusCode:   http.StatusNoContent,
		Status:       "No content",
	}
)

func TestDelete(t *testing.T) {
	subtests := []struct {
		name             string
		newReq           func(verb httpVerb, id uuid.UUID, version *int64) *http.Request
		handleRes        func(response *http.Response, verb httpVerb) (*AccountApiResponse, error)
		apiCall          func(req *http.Request) (*http.Response, error)
		expectedResponse *AccountApiResponse
		expectedErr      error
		expPanic         bool
	}{
		{
			name: "Successfully deleted",
			newReq: func(verb httpVerb, id uuid.UUID, version *int64) *http.Request {
				return &http.Request{
					Header: make(http.Header),
				}
			},
			handleRes: func(response *http.Response, verb httpVerb) (*AccountApiResponse, error) {
				return exp_res_deleted_success, nil
			},
			apiCall: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     "No content",
					StatusCode: http.StatusNoContent,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil
			},
			expectedResponse: exp_res_deleted_success,
			expectedErr:      nil,
		},
		{
			name: "Handle response fails",
			newReq: func(verb httpVerb, id uuid.UUID, version *int64) *http.Request {
				return &http.Request{
					Header: make(http.Header),
				}
			},
			handleRes: func(response *http.Response, verb httpVerb) (*AccountApiResponse, error) {
				return nil, &ApiError{
					StatusCode: 404,
					Status:     "Not found",
				}
			},
			apiCall: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     "Not found",
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil
			},
			expectedResponse: nil,
			expectedErr: &ApiError{
				StatusCode: 404,
				Status:     "Not found",
			},
		},
		{
			name: "Api Call returns error",
			apiCall: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("Failed to do api call")
			},
			newReq: func(verb httpVerb, id uuid.UUID, version *int64) *http.Request {
				return &http.Request{
					Header: make(http.Header),
				}
			},
			expPanic: true,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			newReq = subtest.newReq
			handleRes = subtest.handleRes
			apiCall = subtest.apiCall
			if subtest.expPanic {
				assertPanicDel(t, Delete)
			} else {
				result, err := Delete(uuid.New(), 0)
				if !errors.Is(err, subtest.expectedErr) {
					t.Errorf("expected error (%v), got error (%v)", subtest.expectedErr, err)
				}
				if !reflect.DeepEqual(result, subtest.expectedResponse) {
					t.Errorf("expected (%+v), got (%+v)", subtest.expectedResponse, result)
				}
			}
		})
	}
}

func assertPanicDel(t *testing.T, f func(id uuid.UUID, version int64) (*AccountApiResponse, error)) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f(uuid.New(), 0)
}
