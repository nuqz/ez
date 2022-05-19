package ez

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func RawHTTPRequest(
	client *http.Client,
	method, url string,
	body io.Reader,
) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "can't create new HTTP request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't perform HTTP request")
	}

	return resp, nil
}

type HTTPResponse[T any] struct {
	Raw     *http.Response
	Decoded *T
	decoded bool
}

func NewHTTPResponse[T any](raw *http.Response) *HTTPResponse[T] {
	return &HTTPResponse[T]{Raw: raw}
}

func (response *HTTPResponse[T]) ExpectStatus(status int) error {
	if response.Raw.StatusCode != status {
		return errors.Errorf("expected status %d, but got %d",
			status, response.Raw.StatusCode)
	}

	return nil
}

func HTTPRequest[T any](
	client *http.Client,
	method, url string,
	body io.Reader,
) (*HTTPResponse[T], error) {
	rawResp, err := RawHTTPRequest(client, method, url, body)
	if err != nil {
		return nil, err
	}

	return NewHTTPResponse[T](rawResp), nil
}

type HTTPResponseDecoder[T any] func(io.ReadCloser) (*T, error)

func (response *HTTPResponse[T]) Decode(decoder HTTPResponseDecoder[T]) error {
	if response.decoded {
		return errors.New("already decoded")
	}

	if response.Raw.Body == nil {
		return errors.New("response body is empty")
	}

	defer func() {
		if err := response.Raw.Body.Close(); err != nil {
			// TODO: log an error
			_ = err
		}
	}()

	decoded, err := decoder(response.Raw.Body)
	if err != nil {
		return errors.Wrap(err, "can't decode response")
	}
	response.Decoded = decoded
	response.decoded = true

	return nil
}

type JSONHTTPTransport struct {
	http.RoundTripper
}

func (rt *JSONHTTPTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Accept", "application/json")

	if r.Body != nil {
		r.Header.Set("Content-Type", "application/json")
	}

	return rt.RoundTripper.RoundTrip(r)
}

var (
	JSONHTTPClient = &http.Client{
		Transport: &JSONHTTPTransport{
			RoundTripper: http.DefaultTransport,
		},
	}
)

func JSONRequest[T any](method, url string, body io.Reader) (*T, error) {
	resp, err := HTTPRequest[T](JSONHTTPClient, method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "cant request JSON")
	}

	if err := resp.Decode(FromJSONReader[T]); err != nil {
		return nil, errors.Wrap(err, "cant decode JSON")
	}

	return resp.Decoded, nil
}

func GetJSON[Response any](path string) (*Response, error) {
	return JSONRequest[Response](http.MethodGet, path, nil)
}

func SendJSON[Response any](method, path string, req any) (*Response, error) {
	reqBody, err := ToJSONReader(req)
	if err != nil {
		return nil, err
	}

	return JSONRequest[Response](method, path, reqBody)
}

func PostJSON[Response any](path string, req any) (*Response, error) {
	return SendJSON[Response](http.MethodPost, path, req)
}

func PutJSON[Response any](path string, req any) (*Response, error) {
	return SendJSON[Response](http.MethodPut, path, req)
}
