package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ehox/form3/account"
)

const DIR = "./account/data/"

func main() {
	ac := readFile("sample-account.json")
	apiResponse, err := account.CreateAccount(ac)
	account.Check(err)
	apiRespJson, _ := json.MarshalIndent(apiResponse, "", "  ")
	fmt.Printf("%s\n", apiRespJson)
}

func readFile(name string) account.Account {
	var pathBuilder strings.Builder
	pathBuilder.WriteString(DIR)
	pathBuilder.WriteString(name)
	f, err := os.ReadFile(pathBuilder.String())
	check(err)
	var a account.Account
	err = json.Unmarshal(f, &a)
	check(err)
	defer os.Stdin.Close()
	return a
}

// func newBoolPointer(value bool) *bool {
// 	b := value
// 	return &b
// }

// func newStringPointer(value string) *string {
// 	s := value
// 	return &s
// }

func check(err error) {
	if err != nil {
		panic(err)
	}
}
