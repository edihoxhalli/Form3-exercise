package account

import (
	"github.com/google/uuid"
)

// Fetch enables to get/retrieve an Account record on the form3 API.
// It takes as parameter the id of the record (valid uuid).
// It returns an AccountApiResponse pointer var with the retrieved Account as ResponseBody
// along with the Status and Status Code response details.
// In case any error occurs while attempting to fetch the Account,
// it returns nil, along with the error.
func Fetch(id uuid.UUID) (*AccountApiResponse, error) {
	req, err := newReq(fetchMethod, id, nil)
	if err != nil {
		return nil, err
	}
	response, err := apiCall(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return handleRes(response, fetchMethod)
}
