package account

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func CreateAccount(acc Account) Account {
	accJson, err := json.Marshal(acc)
	check(err)
	requestBody := bytes.NewBuffer(accJson)

	resp, err := http.Post(Url, "application/json", requestBody)
	check(err)
	defer resp.Body.Close()
	var newAccByte []byte
	newAccByte, err = ioutil.ReadAll(resp.Body)
	check(err)
	var newAccStruct Account
	check(json.Unmarshal(newAccByte, &newAccStruct))
	return newAccStruct
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
