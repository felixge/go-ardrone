// Package telnet provides telnet access to the drone.
package telnet

import (
	"fmt"
	gotelnet "github.com/ziutek/telnet"
	"io"
	"io/ioutil"
	"strings"
)

// Dial creates a new telnet connection. The prompt must be set to the shell
// prompt string.
func Dial(addr, prompt string) (*Conn, error) {
	conn, err := gotelnet.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	telnet := &Conn{conn, prompt}
	if err := telnet.waitPrompt(); err != nil {
		return nil, err
	}
	return telnet, err
}

// Conn is a telnet connection.
type Conn struct {
	conn   *gotelnet.Conn
	prompt string
}

// Exec executes the given command and writes the output into the writer.
func (c *Conn) Exec(cmd string, output io.Writer) error {
	if _, err := fmt.Fprintf(c.conn, cmd+"\n"); err != nil {
		return err
	}
	if err := c.discardUntil(cmd + "\r\n"); err != nil {
		return err
	}
	return c.copyUntil(output, c.prompt)
}

// discardUntil discards the output until deli, including the delim itself.
func (c *Conn) discardUntil(delim string) error {
	return c.copyUntil(ioutil.Discard, delim)
}

// discardUntil copies the output until delim, excluding the delim itself.
func (c *Conn) copyUntil(dst io.Writer, delim string) error {
	buf := ""
	for {
		b, err := c.conn.ReadByte()
		if err != nil {
			return err
		}
		buf += string(b)
		if strings.HasSuffix(buf, delim) {
			return nil
		}
		if len(buf) >= len(delim) {
			if _, err := dst.Write([]byte{buf[len(buf)-len(delim)]}); err != nil {
				return err
			}
		}
	}
}

// waitPrompt waits for the prompt to show.
func (c *Conn) waitPrompt() error {
	return c.discardUntil(c.prompt)
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}
