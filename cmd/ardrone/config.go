package main

import (
	"fmt"
	"os"
)

func NewConfigCommand(conf fmt.Stringer) (*ConfigCommand, error) {
	return &ConfigCommand{conf}, nil
}

type ConfigCommand struct {
	conf fmt.Stringer
}

func (c *ConfigCommand) Run() error {
	fmt.Fprintf(os.Stdout, "%s\n", c.conf)
	return nil
}
