package account

import "github.com/google/uuid"

func Delete(id uuid.UUID, version int64) (*AccountApiResponse, error) {
	req := newReq(deleteVerb, id, &version)
	response, err := apiCall(req)

	check(err)
	defer response.Body.Close()
	return handleRes(response, deleteVerb)
}
