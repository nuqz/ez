package http_client

import (
	"io"
	"net/http"
	"testing"

	ezioutils "github.com/nuqz/ez/ioutils"
	ezjson "github.com/nuqz/ez/json"
)

const (
	testServiceURL = "https://api64.ipify.org?format=json"
)

var (
	testJSONClientHeader = map[string]string{
		"User-Agent": "golang-test-suite",
	}
	testJSONClient = &http.Client{
		Transport: NewSingleContentTypeTransport(
			NewConstantHeaderTransport(
				http.DefaultTransport,
				testJSONClientHeader,
			),
			"application/json",
		),
	}
)

type testData struct {
	IP string `json:"ip"`
}

func testRequest(method, url string, body any) (*http.Response, error) {
	var (
		reqBody io.Reader
		err     error
	)

	if method != http.MethodGet && body != nil {
		if reqBody, err = ezjson.ToBuffer(body); err != nil {
			return nil, err
		}
	}

	resp, err := Request(
		testJSONClient,
		http.MethodGet,
		testServiceURL,
		reqBody,
		nil,
	)

	if err != nil {
		return nil, err
	}

	if err := ExpectStatus(http.StatusOK, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func testRequestAndDecode[T any](method, url string, body any) (*T, error) {
	resp, err := testRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return ezioutils.Decode(resp.Body, ezjson.FromReader[T])
}

func TestRequest(t *testing.T) {
	data, err := testRequestAndDecode[testData](http.MethodGet, testServiceURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", data.IP)
}
