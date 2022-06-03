package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

var client *http.Client

func init() {
	client = http.DefaultClient
}

func CreateAccount(acc Account) (*AccountApiResponse, error) {
	accountJSON, err := json.Marshal(acc)
	Check(err)

	request := Client(CREATE)
	request.Header.Add("Content-Type", "application/vnd.api+json")
	request.Header.Add("Content-Length", strconv.Itoa(len([]byte(accountJSON))))
	request.Body = ioutil.NopCloser(bytes.NewReader(accountJSON))
	response, err := client.Do(request)
	Check(err)
	defer response.Body.Close()

	var responseWrapper AccountApiResponse
	responseWrapper.Status = &response.Status
	responseWrapper.StatusCode = &response.StatusCode

	responseBody, err := ioutil.ReadAll(response.Body)
	Check(err)
	if response.StatusCode < 400 {
		var createdAccount Account
		Check(json.Unmarshal(responseBody, &createdAccount))
		responseWrapper.ResponseBody = &createdAccount
		return &responseWrapper, nil
	} else {
		return nil, errors.New(string(responseBody))
	}
}
