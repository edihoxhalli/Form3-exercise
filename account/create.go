package account

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func CreateAccount(acc Account) AccountApiResponse {
	accJson, err := json.Marshal(acc)
	check(err)
	requestBody := bytes.NewBuffer(accJson)

	response, err := http.Post(Url, "application/json", requestBody)
	check(err)

	defer response.Body.Close()
	var responseWrapper AccountApiResponse
	responseWrapper.Status = &response.Status
	responseWrapper.StatusCode = &response.StatusCode
	if response.StatusCode < 400 {
		var newAccByte []byte
		newAccByte, err = ioutil.ReadAll(response.Body)
		check(err)
		var newAccStruct Account
		check(json.Unmarshal(newAccByte, &newAccStruct))
		responseWrapper.ResponseBody = &newAccStruct
	}
	return responseWrapper
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
