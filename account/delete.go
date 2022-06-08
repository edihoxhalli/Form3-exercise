package account

import "github.com/google/uuid"

func Delete(id uuid.UUID, version int64) (*AccountApiResponse, error) {
	req := newRequestWithHeaders(deleteVerb, id, &version)
	response, err := ApiClient.Do(req)

	check(err)
	defer response.Body.Close()
	return handleResponse(response, deleteVerb)
}
