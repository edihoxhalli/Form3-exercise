package account

import (
	"bytes"
	"io/ioutil"
	"strconv"

	"github.com/google/uuid"
)

func Create(acc Account) (*AccountApiResponse, error) {
	accountJSON, err := jsonMarshal(acc)
	check(err)

	request := newReq(createVerb, uuid.Nil, nil)
	request.Header.Add("Content-Type", "application/vnd.api+json")
	request.Header.Add("Content-Length", strconv.Itoa(len([]byte(accountJSON))))

	request.Body = ioutil.NopCloser(bytes.NewReader(accountJSON))
	response, err := apiCall(request)
	check(err)
	defer response.Body.Close()
	return handleRes(response, createVerb)
}
