package types

import (
	"fmt"
)

type HandledError struct {
	Message    string
	InnerError error
}

// Allow HandledError to satisfy error interface.
func (err HandledError) Error() string {
	return err.Message
}
