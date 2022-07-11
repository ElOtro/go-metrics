package repo

import (
	"errors"

	"github.com/ElOtro/go-metrics/internal/repo/postgres"

	"github.com/ElOtro/go-metrics/internal/repo/memory"
)

type Options struct {
	Environment string
}

var ErrEmptyOptions = errors.New("empty options")
var ErrInvalidOptions = errors.New("invalid options")

type Getter interface {
	Get() string
	Set(t, n, v string) error
}

func NewGetter(opts *Options) (Getter, error) {
	if opts == nil {
		return nil, ErrEmptyOptions
	}

	switch opts.Environment {
	case "debug":
		return memory.New(), nil
	case "release":
		return postgres.New(), nil
	default:
		return nil, errors.New("invalid settings")
	}

}
