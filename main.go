package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ehox/form3/account"
	"github.com/google/uuid"
)

const DIR = "./account/data/"

func main() {
	ac := readFile("sample-account.json")
	var ids []string
	for i := 1; i < 4; i++ {
		ac.Data.ID = uuid.New().String()
		apiResponse, err := account.CreateAccount(ac)
		ids = append(ids, ac.Data.ID)

		account.Check(err)
		apiRespJson, _ := json.MarshalIndent(apiResponse, "", "  ")
		fmt.Printf("%s\n%s\n", "Created the following account: ", apiRespJson)
	}

	for _, id := range ids {
		apiResponse, err := account.Fetch(uuid.MustParse(id))
		account.Check(err)
		apiRespJson, _ := json.MarshalIndent(apiResponse, "", "  ")
		fmt.Printf("%s\n%s\n", "Fetched the following account: ", apiRespJson)
	}
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
