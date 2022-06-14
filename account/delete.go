package account

import "github.com/google/uuid"

func Delete(id uuid.UUID, version int64) (*AccountApiResponse, error) {
	req, err := newReq(deleteMethod, id, &version)
	if err != nil {
		return nil, err
	}
	response, err := apiCall(req)
	if err != nil {
		return nil, err
	}
	return handleRes(response, deleteMethod)
}
