// Package account provides a library that can be used as Client of the Form3 API for the resource of Organisation Accounts.
// Current implementation offers Create, Fetch and Delete operations only.
package account

import (
	"bytes"
	"io/ioutil"
	"strconv"

	"github.com/google/uuid"
)

// Create enables to create an Account record on the form3 API.
// It takes as parameter the Account struct the caller wants to create and returns
// the created Account wrapped inside the AccountApiResponse pointer var,
// along with the Status and Status Code response details.
// In case any error occurs while attempting to create the Account,
// it returns nil, along with the error.
func Create(acc Account) (*AccountApiResponse, error) {
	accountJSON, err := jsonMarshal(acc)
	if err != nil {
		return nil, err
	}

	request, err := newReq(createMethod, uuid.Nil, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/vnd.api+json")
	request.Header.Add("Content-Length", strconv.Itoa(len([]byte(accountJSON))))

	request.Body = ioutil.NopCloser(bytes.NewReader(accountJSON))
	response, err := apiCall(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return handleRes(response, createMethod)
}
