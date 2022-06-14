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
	exp_res_fetch_success = &AccountApiResponse{
		ResponseBody: &test_acc,
		StatusCode:   http.StatusOK,
		Status:       "OK",
	}
)

func TestFetch(t *testing.T) {
	subtests := []struct {
		name             string
		newReq           func(verb httpMethod, id uuid.UUID, version *int64) (*http.Request, error)
		handleRes        func(response *http.Response, verb httpMethod) (*AccountApiResponse, error)
		apiCall          func(req *http.Request) (*http.Response, error)
		expectedResponse *AccountApiResponse
		expectedErr      error
	}{
		{
			name: "Successfully fetched",
			newReq: func(verb httpMethod, id uuid.UUID, version *int64) (*http.Request, error) {
				return &http.Request{
					Header: make(http.Header),
				}, nil
			},
			handleRes: func(response *http.Response, verb httpMethod) (*AccountApiResponse, error) {
				return exp_res_fetch_success, nil
			},
			apiCall: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     "No content",
					StatusCode: http.StatusNoContent,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil
			},
			expectedResponse: exp_res_fetch_success,
			expectedErr:      nil,
		},
		{
			name: "Handle response fails",
			newReq: func(verb httpMethod, id uuid.UUID, version *int64) (*http.Request, error) {
				return &http.Request{
					Header: make(http.Header),
				}, nil
			},
			handleRes: func(response *http.Response, verb httpMethod) (*AccountApiResponse, error) {
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
			newReq: func(verb httpMethod, id uuid.UUID, version *int64) (*http.Request, error) {
				return &http.Request{
					Header: make(http.Header),
				}, nil
			},
			expectedErr: errors.New("Failed to do api call"),
		},
		{
			name: "New Request With Headers returns error",
			newReq: func(verb httpMethod, id uuid.UUID, version *int64) (*http.Request, error) {
				return nil, errors.New("Failed to create new request")
			},
			expectedErr: errors.New("Failed to create new request"),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			newReq = subtest.newReq
			handleRes = subtest.handleRes
			apiCall = subtest.apiCall
			result, err := Fetch(uuid.New())
			if err != nil && subtest.expectedErr.Error() != err.Error() {
				t.Errorf("expected error (%v), got error (%v)", subtest.expectedErr, err)
			}
			if !reflect.DeepEqual(result, subtest.expectedResponse) {
				t.Errorf("expected (%+v), got (%+v)", subtest.expectedResponse, result)
			}
		})
	}
}
