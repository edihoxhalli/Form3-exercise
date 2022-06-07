package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ehox/form3/account"
	"github.com/google/uuid"
)

const DIR = "./account/testdata/"

func main() {
	ac := readFile("sample-account.json")
	// ids := []string{
	// "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
	// "f6cea9b7-8fcf-43aa-bd65-ffbb5c219fbb",
	// "c7295128-bffc-4107-b2a9-119bdb407499",
	// }
	// for _, id := range ids {
	// 	apiResponse, err := account.DeleteAccount(uuid.MustParse(id), int64(0))
	// 	panicError(err)
	// 	apiRespJson, _ := json.MarshalIndent(apiResponse, "", "  ")
	// 	fmt.Printf("%s\n%s\n", "Deleted account: ", apiRespJson)
	// }
	var ids []string
	for i := 1; i < 2; i++ {
		ac.Data.ID = uuid.New().String()
		apiResponse, err := account.Create(ac)
		ids = append(ids, ac.Data.ID)

		panicError(err)
		apiRespJson, _ := json.MarshalIndent(apiResponse, "", "  ")
		fmt.Printf("%s\n%s\n", "Created the following account: ", apiRespJson)
	}

	for _, id := range ids {
		apiResponse, err := account.Fetch(uuid.MustParse(id))
		panicError(err)
		apiRespJson, _ := json.MarshalIndent(apiResponse, "", "  ")
		fmt.Printf("%s\n%s\n", "Fetched the following account: ", apiRespJson)
	}
}

func readFile(name string) account.Account {
	var pathBuilder strings.Builder
	pathBuilder.WriteString(DIR)
	pathBuilder.WriteString(name)
	f, err := os.ReadFile(pathBuilder.String())
	panicError(err)
	var a account.Account
	err = json.Unmarshal(f, &a)
	panicError(err)
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

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
