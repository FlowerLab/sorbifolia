package errors

import (
	"errors"
	"strings"
)

type Error []error

func (e Error) Error() string {
	s := new(strings.Builder)
	s.WriteString("errors:")
	for _, err := range e {
		s.WriteString("\n\t")
		s.WriteString(err.Error())
	}
	return s.String()
}

func (e Error) Errors() []error { return e }

func (e Error) Is(target error) bool {
	for _, err := range e {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e Error) As(target any) bool {
	for _, err := range e {
		if ok := errors.As(err, target); ok {
			return true
		}
	}
	return false
}

func (e Error) Index(target error) int {
	for idx, err := range e {
		if errors.Is(err, target) {
			return idx
		}
	}
	return -1
}
