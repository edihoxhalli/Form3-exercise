package account

import (
	"github.com/google/uuid"
)

func Fetch(id uuid.UUID) (*AccountApiResponse, error) {
	req, err := newReq(fetchMethod, id, nil)
	if err != nil {
		return nil, err
	}
	response, err := apiCall(req)
	if err != nil {
		return nil, err
	}
	return handleRes(response, fetchMethod)
}
