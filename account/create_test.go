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
	uk_account_without_cop = Account{
		Data: &AccountData{
			Attributes: &AccountAttributes{
				Country:             "GB",
				BaseCurrency:        "GBP",
				BankID:              "400300",
				BankIDCode:          "GBDSC",
				Bic:                 "NWBKGB22",
				ValidationType:      "card",
				ReferenceMask:       "############",
				AcceptanceQualifier: "same_day",
				UserDefinedData: &[]UserDefinedData{
					{
						Key:   "Some account related key",
						Value: "Some account related value",
					},
				},
			},
			Type:           "accounts",
			ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		},
	}
	exp_res_success = &AccountApiResponse{
		ResponseBody: &uk_account_without_cop,
		StatusCode:   http.StatusCreated,
		Status:       "Created",
	}
)

func TestCreate(t *testing.T) {

	subtests := []struct {
		name             string
		newReq           func(verb httpVerb, id uuid.UUID, version *int64) *http.Request
		handleRes        func(response *http.Response, verb httpVerb) (*AccountApiResponse, error)
		apiCall          func(req *http.Request) (*http.Response, error)
		jsonMarshal      func(v any) ([]byte, error)
		expectedResponse *AccountApiResponse
		expectedErr      error
		expPanic         bool
	}{
		{
			name: "Successfully created",
			newReq: func(verb httpVerb, id uuid.UUID, version *int64) *http.Request {
				return &http.Request{
					Header: make(http.Header),
				}
			},
			handleRes: func(response *http.Response, verb httpVerb) (*AccountApiResponse, error) {
				return exp_res_success, nil
			},
			apiCall: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     "Created",
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewBufferString("body")),
				}, nil
			},
			jsonMarshal:      json.Marshal,
			expectedResponse: exp_res_success,
			expectedErr:      nil,
		},
		{
			name: "Json Marshaling returns error",
			jsonMarshal: func(v any) ([]byte, error) {
				return []byte("Ignore"), errors.New("Failed to marshall")
			},
			expPanic: true,
		},
		{
			name:        "Api Call returns error",
			jsonMarshal: json.Marshal,
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
			jsonMarshal = subtest.jsonMarshal
			if subtest.expPanic {
				assertPanic(t, Create, uk_account_without_cop)
			} else {
				result, err := Create(uk_account_without_cop)
				if err != nil {
					t.Errorf("expected error (%v), got error (%v)", subtest.expectedErr, err)
				}
				if !reflect.DeepEqual(result, subtest.expectedResponse) {
					t.Errorf("expected (%+v), got (%+v)", subtest.expectedResponse, result)
				}
			}
		})
	}
}

func assertPanic(t *testing.T, f func(acc Account) (*AccountApiResponse, error), acc Account) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f(acc)
}
