package http_client

import "net/http"

type SingleContentTypeTransport struct {
	http.RoundTripper
	contentType string
}

func NewSingleContentTypeTransport(
	rt http.RoundTripper,
	contentType string,
) *SingleContentTypeTransport {
	return &SingleContentTypeTransport{
		RoundTripper: rt,
		contentType:  contentType,
	}
}

func (t *SingleContentTypeTransport) RoundTrip(
	r *http.Request,
) (*http.Response, error) {
	r.Header.Set("Accept", t.contentType)

	if r.Body != nil {
		r.Header.Set("Content-Type", t.contentType)
	}

	return t.RoundTripper.RoundTrip(r)
}

type ConstantHeaderTransport struct {
	http.RoundTripper
	header map[string]string
}

func NewConstantHeaderTransport(
	rt http.RoundTripper,
	headers map[string]string,
) *ConstantHeaderTransport {
	return &ConstantHeaderTransport{
		RoundTripper: rt,
		header:       headers,
	}
}

func (t *ConstantHeaderTransport) RoundTrip(
	r *http.Request,
) (*http.Response, error) {
	for name, value := range t.header {
		r.Header.Set(name, value)
	}

	return t.RoundTripper.RoundTrip(r)
}
