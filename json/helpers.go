package json

import (
	"bytes"
	"encoding/json"
	"io"
)

func FromBytes[T any](body []byte) (*T, error) {
	out := new(T)

	if err := json.Unmarshal(body, out); err != nil {
		return nil, err
	}

	return out, nil
}

func FromReader[T any](body io.ReadCloser) (*T, error) {
	defer body.Close()
	out := new(T)

	if err := json.NewDecoder(body).Decode(out); err != nil {
		return nil, err
	}

	return out, nil
}

func ToBuffer(body any) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	return buf, nil
}
