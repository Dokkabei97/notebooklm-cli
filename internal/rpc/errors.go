package rpc

import "fmt"

type ErrorCode int

const (
	ErrUnknown       ErrorCode = 0
	ErrAuth          ErrorCode = 401
	ErrForbidden     ErrorCode = 403
	ErrNotFound      ErrorCode = 404
	ErrRateLimit     ErrorCode = 429
	ErrServer        ErrorCode = 500
	ErrUnavailable   ErrorCode = 503
)

type Error struct {
	Code    ErrorCode
	Message string
	Method  string
}

func (e *Error) Error() string {
	if e.Method != "" {
		return fmt.Sprintf("rpc %s: [%d] %s", e.Method, e.Code, e.Message)
	}
	return fmt.Sprintf("rpc error [%d]: %s", e.Code, e.Message)
}

func IsAuthError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrAuth || e.Code == ErrForbidden
	}
	return false
}

func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrNotFound
	}
	return false
}

func IsRateLimit(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrRateLimit
	}
	return false
}
