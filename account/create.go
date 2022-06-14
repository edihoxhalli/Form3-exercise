package account

import (
	"bytes"
	"io/ioutil"
	"strconv"

	"github.com/google/uuid"
)

func Create(acc Account) (*AccountApiResponse, error) {
	accountJSON, err := jsonMarshal(acc)
	if err != nil {
		return nil, err
	}

	request, err := newReq(createMethod, uuid.Nil, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/vnd.api+json")
	request.Header.Add("Content-Length", strconv.Itoa(len([]byte(accountJSON))))

	request.Body = ioutil.NopCloser(bytes.NewReader(accountJSON))
	response, err := apiCall(request)
	if err != nil {
		return nil, err
	}
	return handleRes(response, createMethod)
}
