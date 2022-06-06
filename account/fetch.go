package account

import (
	"github.com/google/uuid"
)

func FetchAccount(id uuid.UUID) (*AccountApiResponse, error) {
	req := newRequestWithHeaders(FETCH, id, nil)
	response, err := ApiClient.Do(req)

	check(err)
	defer response.Body.Close()
	return handleResponse(response, FETCH)
}
