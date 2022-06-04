package account

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/google/uuid"
)

func CreateAccount(acc Account) (*AccountApiResponse, error) {
	accountJSON, err := json.Marshal(acc)
	Check(err)

	request := NewRequestWithHeaders(CREATE, uuid.Nil)
	request.Header.Add("Content-Type", "application/vnd.api+json")
	request.Header.Add("Content-Length", strconv.Itoa(len([]byte(accountJSON))))

	request.Body = ioutil.NopCloser(bytes.NewReader(accountJSON))
	response, err := ApiClient.Do(request)
	Check(err)
	defer response.Body.Close()
	return HandleResponse(response, CREATE)
}
