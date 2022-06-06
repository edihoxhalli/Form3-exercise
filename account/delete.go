package account

import "github.com/google/uuid"

func DeleteAccount(id uuid.UUID, version int64) (*AccountApiResponse, error) {
	req := newRequestWithHeaders(DELETE, id, &version)
	response, err := ApiClient.Do(req)

	check(err)
	defer response.Body.Close()
	return handleResponse(response, DELETE)
}
