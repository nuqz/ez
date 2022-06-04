package ioutils

import (
	"io"

	"github.com/pkg/errors"
)

var (
	ErrNoSource = errors.New("please, provide source to decode")
)

type Decoder[T any] func(io.ReadCloser) (*T, error)

func Decode[T any](source io.ReadCloser, decoder Decoder[T]) (*T, error) {
	if source == nil {
		return nil, ErrNoSource
	}

	defer func() {
		if err := source.Close(); err != nil {
			// TODO: log an error
			_ = err
		}
	}()

	decoded, err := decoder(source)
	if err != nil {
		return nil, errors.Wrap(err, "can't decode response")
	}

	return decoded, nil
}
