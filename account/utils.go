package account

import (
	"net/http"
	"time"
)

const Url string = "http://localhost:8080/v1/organisation/accounts"

type Verb int

const (
	CREATE Verb = iota
	FETCH
	DELETE
)

func (index Verb) String() string {
	return [...]string{"POST", "GET", "DELETE"}[index]
}

func Client(verb Verb) *http.Request {
	req, err := http.NewRequest(verb.String(), Url, nil)
	Check(err)

	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Date", time.Now().String())
	req.Header.Add("Accept", "application/vnd.api+json")
	return req
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
