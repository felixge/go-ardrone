package main

import (
	"fmt"
	"github.com/felixge/go-ardrone/net/ftp"
	"github.com/felixge/go-ardrone/net/telnet"
	"os"
	"os/exec"
	"path/filepath"
)

// TestCommandConfig defines the config required for the test command.
type TestCommandConfig interface {
	DialFtp() (*ftp.Conn, error)
	DialTelnet() (*telnet.Conn, error)
	Wd() string
}

// NewTestCommand returns a new TestCommand.
func NewTestCommand(args []string, config TestCommandConfig) (*TestCommand, error) {
	return &TestCommand{config: config}, nil
}

// TestCommand cross compiles the go tests of a pkg and runs them on the drone.
type TestCommand struct {
	config TestCommandConfig
}

// Run implements the Command interface.
func (c *TestCommand) Run() error {
	localPath, err := c.build()
	if err != nil {
		return err
	}

	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	ftpConn, err := c.config.DialFtp()
	if err != nil {
		return err
	}
	defer ftpConn.Close()

	_, name := filepath.Split(localPath)
	dronePath, err := ftpConn.Upload(file, name)
	if err != nil {
		return err
	}
	defer ftpConn.Del(name)

	telnetConn, err := c.config.DialTelnet()
	if err != nil {
		return err
	}
	defer telnetConn.Close()

	cmd := fmt.Sprintf("chmod +x %s && %s", dronePath, dronePath)
	return telnetConn.Exec(cmd, os.Stdout)
}

// build cross compiles the test binary and returns the path to it.
func (c *TestCommand) build() (string, error) {
	path := "go" // go must be in $PATH
	// -c compiles without running the test
	args := append([]string{"test", "-c"})
	cmd := exec.Command(path, args...)
	pkgDir := c.config.Wd()
	cmd.Dir = pkgDir
	cmd.Env = CrossEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	_, pkgName := filepath.Split(pkgDir)
	fileName := filepath.Join(pkgDir, pkgName+".test")
	return fileName, err
}
