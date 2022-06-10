package account

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHttpVerbString(t *testing.T) {
	subtests := []struct {
		name      string
		verb      httpVerb
		expString string
	}{
		{
			"Verb POST mapped successfully",
			createVerb,
			"POST",
		},
		{
			"Verb FETCH mapped successfully",
			fetchVerb,
			"GET",
		},
		{
			"Verb DELETE mapped successfully",
			deleteVerb,
			"DELETE",
		},
	}
	for _, subtest := range subtests {
		t.Run("Http Verbs mapped successfully", func(t *testing.T) {
			if subtest.verb.String() != subtest.expString {
				t.Errorf("expected httpVerb %d to produce %s", subtest.verb, subtest.verb.String())
			}
		})
	}
}

func TestNewReq(t *testing.T) {

	subtests := []struct {
		name            string
		httpNewRequest  func(method, url string, body io.Reader) (*http.Request, error)
		httpVerb        httpVerb
		id              uuid.UUID
		version         *int64
		expHostHdr      string
		expAcceptHdr    string
		expVersionParam string
		expPanic        bool
	}{
		{
			name: "New POST request with Headers",
			httpNewRequest: func(method, custUrl string, body io.Reader) (*http.Request, error) {
				return &http.Request{
					Header: make(http.Header),
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme:  "http://",
						Host:    "localhost:8080",
						RawPath: "/accounts",
					}}, nil
			},
			httpVerb:     createVerb,
			id:           uuid.New(),
			expHostHdr:   Host,
			expAcceptHdr: "application/vnd.api+json",
		},
		{
			name: "New DELETE request with Headers",
			httpNewRequest: func(method, custUrl string, body io.Reader) (*http.Request, error) {
				return &http.Request{
					Header: make(http.Header),
					Method: http.MethodDelete,
					URL: &url.URL{
						Scheme:  "http://",
						Host:    "localhost:8080",
						RawPath: "/accounts",
					}}, nil
			},
			httpVerb:        deleteVerb,
			id:              uuid.New(),
			version:         int64ToPointer(0),
			expVersionParam: "0",
			expHostHdr:      Host,
			expAcceptHdr:    "application/vnd.api+json",
		},
		{
			name: "http.NewRequest returns error and func panics",
			httpNewRequest: func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("Error")
			},
			expPanic: true,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			httpNewRequest = subtest.httpNewRequest
			if subtest.expPanic {
				assertPanicNewReq(t, newReq)
			} else {
				result := newReq(subtest.httpVerb, subtest.id, subtest.version)
				if result.Header.Get("Host") != subtest.expHostHdr {
					t.Errorf("expected header Host (%+v), got (%+v)", subtest.expHostHdr, result.Header.Get("Host"))
				}
				if result.Header.Get("Accept") != subtest.expAcceptHdr {
					t.Errorf("expected header Accept (%+v), got (%+v)", subtest.expAcceptHdr, result.Header.Get("Accept"))
				}
				testDateHdr(result.Header.Get("Date"), t)
				versionQueryParam := result.URL.Query().Get("version")
				if versionQueryParam != subtest.expVersionParam {
					t.Errorf("expected version param (%+v), got (%+v)", subtest.expVersionParam, versionQueryParam)
				}
			}
		})
	}
}

func testDateHdr(dateHdrStr string, t *testing.T) {
	// layout := "Mon Jan 2 15:04:05 MST 2006  (MST is GMT-0700)"
	dateHdrTime, err := time.Parse(time.RFC3339Nano, dateHdrStr)
	if err != nil {
		t.Errorf("got error from Date header conversion %s", err.Error())
	}
	timeNow := time.Now()
	diff := timeNow.Sub(dateHdrTime)
	if diff*time.Second < 1 {
		t.Errorf("expected header Date within 1 sec before (%s), got (%s)", timeNow.Format(time.RFC3339Nano), dateHdrStr)
	}
}

func int64ToPointer(i int64) *int64 {
	return &i
}

func assertPanicNewReq(t *testing.T, f func(verb httpVerb, id uuid.UUID, version *int64) *http.Request) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f(createVerb, uuid.Nil, nil)
}
