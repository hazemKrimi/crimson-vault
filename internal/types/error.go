package types

import (
	"fmt"
	"strings"
)

type Error struct {
	Messages []string
	Cause    error
	Code     int
}

func (err Error) Error() string {
	return fmt.Sprintf("%s",
		strings.Join(err.Messages, ", "))
}
