package stringsvc

import (
	"errors"
	"strings"
)

type stringService struct{}

// NewService returns a new userService
func NewService() StringService {
	return &stringService{}
}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

// ErrEmpty is returned when input string is empty
var ErrEmpty = errors.New("Empty string")
