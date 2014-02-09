// Command ardrone implements the ardrone command line utility.
package main

import (
	"fmt"
	"github.com/felixge/go-ardrone/config"
	"os"
)

func Usage() string {
	return "ardrone <cmd>\n\n" +
		"Where <cmd> is one of:\n\n" +
		"  env  Prints the env configuration.\n" +
		"  test Builds go tests in the current directory and runs them the drone."
}

// main executes the program.
func main() {
	conf, err := NewConfigFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config: %s\n", err)
		os.Exit(1)
	}
	cmd, err := NewCommand(os.Args, conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
		fmt.Fprintf(os.Stderr, "Usage: %s\n", Usage())
		os.Exit(2)
	}
	if err = cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(3)
	}
}

// fatal prints the given error and terminates the program.
func fatalUsage(err error) {
}

// Command defines a command of the ardrone utility.
type Command interface {
	// Run executes the command.
	Run() error
}

// NewConfigFromEnv returns the config from the environment.
func NewConfigFromEnv() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	conf, err := config.NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return &Config{conf, wd}, nil
}

// Config defines the ardrone command config.
type Config struct {
	*config.Config
	wd string
}

// Wd returns the working directory.
func (c *Config) Wd() string {
	return c.wd
}

// NewCommand returns the command for the given arguments.
func NewCommand(args []string, conf *Config) (Command, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("No command was given.")
	}
	cmd := args[1]
	args = args[2:]
	switch cmd {
	case "config":
		return NewConfigCommand(conf)
	case "test":
		return NewTestCommand(args, conf)
	default:
		return nil, fmt.Errorf("Unknown command: %s", cmd)
	}
}

// CrossEnv returns environment variables suitable for cross compiling Go code
// to run on the ardrone.
func CrossEnv() []string {
	// @TODO Not sure what happens if os.Environ() already contains one of these
	// names.  I think it depends on what exec(3) does, but the man page didn't
	// seem to tell me.
	return append(os.Environ(), "GOOS=linux", "GOARCH=arm")
}
