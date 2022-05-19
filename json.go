package ez

import (
	"bytes"
	"encoding/json"
	"io"
)

func FromJSONBytes[T any](body []byte) (*T, error) {
	out := new(T)

	if err := json.Unmarshal(body, out); err != nil {
		return nil, err
	}

	return out, nil
}

func FromJSONReader[T any](body io.ReadCloser) (*T, error) {
	defer body.Close()
	out := new(T)

	if err := json.NewDecoder(body).Decode(out); err != nil {
		return nil, err
	}

	return out, nil
}

func ToJSONReader(body any) (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	return buf, nil
}
