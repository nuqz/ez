package http_client

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func Request(
	client *http.Client,
	method, url string,
	body io.Reader,
	requestModifier func(*http.Request) *http.Request,
) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "can't create new HTTP request")
	}

	if requestModifier != nil {
		request = requestModifier(request)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "can't perform HTTP request")
	}

	return response, nil
}

func ExpectStatus(status int, response *http.Response) error {
	if response.StatusCode != status {
		return errors.Errorf("expected status %d, but got %d",
			status, response.StatusCode)
	}
	return nil
}
