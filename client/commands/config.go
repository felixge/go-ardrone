package commands

import (
	"fmt"
)

type Config struct {
	Key string
	Value string
}

func (this Config) String(number int) string {
	return fmt.Sprintf("AT*CONFIG=%d,\"%s\",\"%s\"\r", number, this.Key, this.Value)
}
