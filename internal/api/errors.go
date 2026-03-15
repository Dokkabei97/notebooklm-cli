package api

import "fmt"

type Error struct {
	Op      string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func wrapErr(op, msg string, err error) error {
	return &Error{Op: op, Message: msg, Err: err}
}
