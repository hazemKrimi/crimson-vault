package types

import "fmt"

type Error struct {
	Messages []string
	Code     int
}

func (err Error) Error() string {
	return fmt.Sprintf("%v",
		err.Messages)
}
