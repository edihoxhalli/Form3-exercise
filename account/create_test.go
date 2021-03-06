package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

var (
	test_acc = Account{
		Data: &AccountData{
			Attributes: &AccountAttributes{
				Country:      newStringPointer("GB"),
				BaseCurrency: "GBP",
				BankID:       "400300",
				BankIDCode:   "GBDSC",
				Bic:          "NWBKGB22",
			},
			Type:           "accounts",
			ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		},
	}
	exp_res_created_success = &AccountApiResponse{
		ResponseBody: &test_acc,
		StatusCode:   http.StatusCreated,
		Status:       "Created",
	}
)

func newStringPointer(s string) *string {
	return &s
}

func TestCreate(t *testing.T) {

	subtests := []struct {
		name             string
		newReq           func(verb httpMethod, id uuid.UUID, version *int64) (*http.Request, error)
		handleRes        func(response *http.Response, verb httpMethod) (*AccountApiResponse, error)
		apiCall          func(req *http.Request) (*http.Response, error)
		jsonMarshal      func(v any) ([]byte, error)
		expectedResponse *AccountApiResponse
		expectedErr      error
	}{
		{
			name: "Successfully created",
			newReq: func(verb httpMethod, id uuid.UUID, version *int64) (*http.Request, error) {
				return &http.Request{
					Header: make(http.Header),
				}, nil
			},
			handleRes: func(response *http.Response, verb httpMethod) (*AccountApiResponse, error) {
				return exp_res_created_success, nil
			},
			apiCall: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     "Created",
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewBufferString("body")),
				}, nil
			},
			jsonMarshal:      json.Marshal,
			expectedResponse: exp_res_created_success,
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
					StatusCode: 400,
					Status:     "Bad request",
				}
			},
			apiCall: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     "Bad request",
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil
			},
			jsonMarshal:      json.Marshal,
			expectedResponse: nil,
			expectedErr: &ApiError{
				StatusCode: 400,
				Status:     "Bad request",
			},
		},
		{
			name: "Json Marshaling returns error",
			jsonMarshal: func(v any) ([]byte, error) {
				return nil, errors.New("Failed to marshall")
			},
			expectedErr: errors.New("Failed to marshall"),
		},
		{
			name:        "Api Call returns error",
			jsonMarshal: json.Marshal,
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
			name:        "New Request With Headers returns error",
			jsonMarshal: json.Marshal,
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
			jsonMarshal = subtest.jsonMarshal
			result, err := Create(test_acc)
			if err != nil && subtest.expectedErr.Error() != err.Error() {
				t.Errorf("expected error (%v), got error (%v)", subtest.expectedErr, err)
			}
			if !reflect.DeepEqual(result, subtest.expectedResponse) {
				t.Errorf("expected (%+v), got (%+v)", subtest.expectedResponse, result)
			}
		})
	}
}
