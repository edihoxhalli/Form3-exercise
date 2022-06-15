package account

import "github.com/google/uuid"

// Delete enables to delete an Account record on the form3 API.
// It takes as parameter the id of the record (valid uuid) and the version.
// It returns an AccountApiResponse pointer var with nil as ResponseBody
// along with the Status and Status Code response details.
// In case any error occurs while attempting to delete the Account,
// it returns nil, along with the error.
func Delete(id uuid.UUID, version int64) (*AccountApiResponse, error) {
	req, err := newReq(deleteMethod, id, &version)
	if err != nil {
		return nil, err
	}
	response, err := apiCall(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return handleRes(response, deleteMethod)
}
