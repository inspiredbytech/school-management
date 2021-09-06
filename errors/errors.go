package errors

import "errors"

var (
	ErrInvalidArgument = errors.New("Invalid function argument(s)")
	ErrSchoolNotFound  = errors.New("School not found")
)
