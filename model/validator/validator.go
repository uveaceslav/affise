package validator

import (
	"errors"
	"github.com/uveaceslav/affise/model"
)

const (
	// MaxURLs is the maximum allowable number of URLs in a request,
	MaxURLs = 20 // ToDo: Move to config
)

var (
	// ErrExceededMaxURLs is an error indicating that the request has more
	// than the allowable MaxURLs URL entries.
	ErrExceededMaxURLs = errors.New("Exceeded maximum number of URLs in a request")
)

type RequestEntityValidator interface {
	Validate(re model.RequestEntity) error
}

type requestEntityValidator struct {
}

func NewRequestEntityValidator() RequestEntityValidator {
	return &requestEntityValidator{}
}

func (rev *requestEntityValidator) Validate(re model.RequestEntity) error {
	if len(re.URLs) > MaxURLs {
		return ErrExceededMaxURLs
	}

	return nil
}
