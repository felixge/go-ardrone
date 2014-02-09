// Package ftp provides ftp access to the ardrone.
package ftp

import (
	goftp "github.com/jlaffaye/goftp"
	"io"
	"path/filepath"
)

// Dial connects to the ardrone ftp server at addr. The ftpHome must be set
// to the ftp home directory.
func Dial(addr string, ftpHome string) (*Conn, error) {
	ftpConn, err := goftp.Connect(addr)
	if err != nil {
		return nil, err
	}
	return &Conn{conn: ftpConn, ftpHome: ftpHome}, nil
}

// Conn is a ftp connection.
type Conn struct {
	conn    *goftp.ServerConn
	ftpHome string
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.conn.Quit()
}

// Upload uploads the data to the relative ftp path. Returns the absolute path
// of the file on the drone.
func (c *Conn) Upload(data io.Reader, path string) (string, error) {
	err := c.conn.Stor(path, data)
	return filepath.Join(c.ftpHome, path), err
}

// Del deletes the given relative ftp path.
func (c *Conn) Del(path string) error {
	return c.conn.Delete(path)
}
