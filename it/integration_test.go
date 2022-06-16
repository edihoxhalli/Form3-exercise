// Package it holds Integration (e2e) Tests for package account.
package it

import (
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/edihoxhalli/Form3-exercise/account"
	"github.com/google/uuid"
)

var (
	l        = log.Default()
	test_acc = account.Account{
		Data: &account.AccountData{
			Attributes: &account.AccountAttributes{
				AccountClassification: newStringPointer("Personal"),
				AccountMatchingOptOut: newBoolPointer(false),
				AccountNumber:         "1231555",
				AlternativeNames: []string{
					"Sam Holder",
				},
				BankID:       "400300",
				BankIDCode:   "GBDSC",
				BaseCurrency: "GBP",
				Bic:          "NWBKGB22",
				Country:      newStringPointer("GB"),
				Iban:         "GB11NWBK40030041426819",
				JointAccount: newBoolPointer(false),
				Name: []string{
					"Samantha Holder",
				},
				SecondaryIdentification: "A1B2C3D4",
				Status:                  newStringPointer("pending"),
				Switched:                newBoolPointer(true),
			},
			Type:           "accounts",
			ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		},
	}
	exp_res_created_success = &account.AccountApiResponse{
		ResponseBody: &test_acc,
		StatusCode:   http.StatusCreated,
		Status:       "201 Created",
	}
	exp_res_fetch_success = &account.AccountApiResponse{
		ResponseBody: &test_acc,
		StatusCode:   http.StatusOK,
		Status:       "200 OK",
	}
	exp_res_deleted_success = &account.AccountApiResponse{
		ResponseBody: nil,
		StatusCode:   http.StatusNoContent,
		Status:       "204 No Content",
	}
)

func newBoolPointer(value bool) *bool {
	return &value
}

func newStringPointer(s string) *string {
	return &s
}

func newInt64Pointer(i int64) *int64 {
	return &i
}

func init() {
	account.Host = "http://accountapi:8080/"
}

func TestCreateE2E(t *testing.T) {
	upd := exp_res_created_success.ResponseBody
	upd.Data.Version = newInt64Pointer(0)
	exp_res_created_success.ResponseBody = upd

	res, err := account.Create(test_acc)
	if err != nil {
		t.Errorf("expected nil error, got (%v)", err)
	} else {
		l.Printf("Created account (%v)", spew.Sdump(*res))
	}
	if !reflect.DeepEqual(res, exp_res_created_success) {
		t.Errorf("expected created account response (%+v), got (%+v)", spew.Sdump(exp_res_created_success), spew.Sdump(*res))
	}
}

func TestFetchE2E(t *testing.T) {
	upd := exp_res_fetch_success.ResponseBody
	upd.Data.Version = newInt64Pointer(0)
	exp_res_created_success.ResponseBody = upd

	res, err := account.Fetch(uuid.MustParse(test_acc.Data.ID))
	if err != nil {
		t.Errorf("expected nil error, got (%v)", err)
	} else {
		l.Printf("Fetched account (%v)", spew.Sdump(*res))
	}
	if !reflect.DeepEqual(res, exp_res_fetch_success) {
		t.Errorf("expected fetched account response (%+v), got (%+v)", spew.Sdump(exp_res_fetch_success), spew.Sdump(*res))
	}
}

func TestDeleteE2E(t *testing.T) {
	res, err := account.Delete(uuid.MustParse(test_acc.Data.ID), 0)
	if err != nil {
		t.Errorf("expected nil error, got (%v)", err)
	} else {
		l.Printf("Deleted account (%v)", spew.Sdump(*res))
	}
	if !reflect.DeepEqual(res, exp_res_deleted_success) {
		t.Errorf("expected deleted account response (%+v), got (%+v)", spew.Sdump(exp_res_deleted_success), spew.Sdump(res))
	}
}
