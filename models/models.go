package models

import "fmt"

type HandledError struct {
	StatusCode int
	Message    string
}

func (e HandledError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}
