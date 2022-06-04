package account

import (
	"github.com/google/uuid"
)

func FetchAccount(id uuid.UUID) (*AccountApiResponse, error) {
	req := NewRequestWithHeaders(FETCH, id, nil)
	response, err := ApiClient.Do(req)

	Check(err)
	defer response.Body.Close()
	return HandleResponse(response, FETCH)
}
