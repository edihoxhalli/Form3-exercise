package account

import (
	"github.com/google/uuid"
)

func Fetch(id uuid.UUID) (*AccountApiResponse, error) {
	req := newReq(fetchVerb, id, nil)
	response, err := ApiClient.Do(req)

	check(err)
	defer response.Body.Close()
	return handleRes(response, fetchVerb)
}
