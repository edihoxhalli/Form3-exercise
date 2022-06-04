package account

import (
	"github.com/google/uuid"
)

func Fetch(id uuid.UUID) (*AccountApiResponse, error) {
	req := NewRequestWithHeaders(FETCH, id)
	response, err := ApiClient.Do(req)

	Check(err)
	defer response.Body.Close()
	return HandleResponse(response, FETCH)
}
