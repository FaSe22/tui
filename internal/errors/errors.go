package errors

import "fmt"

type UIError struct {
	Component string
	Message   string
	Err       error
}

func (e *UIError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Component, e.Message, e.Err)
}
