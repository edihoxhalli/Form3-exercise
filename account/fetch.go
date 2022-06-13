package account

import (
	"github.com/google/uuid"
)

func Fetch(id uuid.UUID) (*AccountApiResponse, error) {
	req := newReq(fetchMethod, id, nil)
	response, err := apiCall(req)

	check(err)
	defer response.Body.Close()
	return handleRes(response, fetchMethod)
}
