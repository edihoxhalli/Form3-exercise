package account

import "github.com/google/uuid"

func DeleteAccount(id uuid.UUID, version int64) (*AccountApiResponse, error) {
	req := NewRequestWithHeaders(DELETE, id, &version)
	response, err := ApiClient.Do(req)

	Check(err)
	defer response.Body.Close()
	return HandleResponse(response, DELETE)
}
