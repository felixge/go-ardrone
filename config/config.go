// Package config provides configured instances. Right now all configuration is
// done via environment variables.
package config

import (
	"bytes"
	"fmt"
	"github.com/felixge/go-ardrone/net/ftp"
	"github.com/felixge/go-ardrone/net/telnet"
	"net"
	"os"
	"strings"
)

type envKey string

const (
	host         envKey = "ARDRONE_HOST"
	ftpPort      envKey = "ARDRONE_FTP_PORT"
	ftpHome      envKey = "ARDRONE_FTP_HOME"
	telnetPort   envKey = "ARDRONE_TELNET_PORT"
	telnetPrompt envKey = "ARDRONE_TELNET_PROMPT"
)

// NewConfigFromEnv loads the config from env variables.
func NewConfigFromEnv() (*Config, error) {
	config := DefaultConfig()
	for key, _ := range config.vals {
		envVal := os.Getenv(string(key))
		if envVal != "" {
			config.vals[key] = envVal
		}
	}
	return config, nil
}

// getenv reads the environment variable named key into val if not empty.
func getenv(key string, val *string) {
	envVal := os.Getenv(key)
	if envVal != "" {
		*val = envVal
	}
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		vals: map[envKey]string{
			host:         "192.168.1.1",
			ftpPort:      "21",
			ftpHome:      "/data/video",
			telnetPort:   "23",
			telnetPrompt: "# ",
		},
	}
}

// Config defines the ardrone config.
type Config struct {
	vals map[envKey]string
}

func (c *Config) get(key envKey) string {
	return c.vals[key]
}

// DialFtp creates a new ftp connection.
func (c *Config) DialFtp() (*ftp.Conn, error) {
	addr := net.JoinHostPort(c.get(host), c.get(ftpPort))
	return ftp.Dial(addr, c.get(ftpHome))
}

// DialTelnet creates a new telnet connection.
func (c *Config) DialTelnet() (*telnet.Conn, error) {
	addr := net.JoinHostPort(c.get(host), c.get(telnetPort))
	return telnet.Dial(addr, c.get(telnetPrompt))
}

// Strings prints a human readable version of the config.
func (c *Config) String() string {
	buf := &bytes.Buffer{}
	for key, val := range c.vals {
		fmt.Fprintf(buf, "%s=%s\n", key, val)
	}
	return strings.TrimSpace(buf.String())
}
