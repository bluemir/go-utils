package auth

import (
	"fmt"

	"github.com/bluemir/go-utils/auth/codes"
)

type AuthError struct {
	error
	code codes.ErrorCode
}

func Errorf(code codes.ErrorCode, format string, val ...interface{}) error {
	return Error(code, fmt.Errorf(format, val...))
}
func Error(code codes.ErrorCode, err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*AuthError); ok {
		return e
	}
	return &AuthError{err, code}
}

func ErrorCode(err error) codes.ErrorCode {
	if err == nil {
		return codes.None
	}
	if e, ok := err.(*AuthError); ok {
		return e.code
	}
	return codes.Unknown
}
