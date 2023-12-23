package controller

import (
	"fmt"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Controller struct {
	*textproto.Conn

	mu sync.Mutex

	lastSignal     Signal
	lastSignalTime time.Time
}

// NewController creates a new controller connection to the given address.
// The address should be in the form "host:port". If the address is empty,
// the default address will be 127.0.0.1:9051.
func NewController(addr string) (*Controller, error) {
	if addr == "" {
		addr = "127.0.0.1:9051"
	}
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Controller{Conn: conn, lastSignal: -1}, nil
}

// Signal sends a signal to the Tor process.
// If the same signal is sent twice within 10 seconds, the second signal is
// ignored.
// See https://gitweb.torproject.org/torspec.git/tree/control-spec.txt#n102
func (c *Controller) Signal(signal Signal) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if signal == c.lastSignal && time.Since(c.lastSignalTime) < 10*time.Second+500*time.Millisecond {
		return nil // because this is a no-op
	}
	_, _, err := c.makeRequest("SIGNAL " + signal.String())
	if err != nil {
		return err
	}
	c.lastSignal = signal
	c.lastSignalTime = time.Now()
	return nil
}

func (c *Controller) GetAddress() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.getInfo("address")
}

func (c *Controller) GetBytesRead() (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.getInfoInt("traffic/read")
}

func (c *Controller) GetBytesWritten() (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.getInfoInt("traffic/written")
}

func (c *Controller) GetVersion() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.getInfo("version")
}

// Authenticate authenticates the controller connection
// If the password is empty, it will authenticate without a password.
func (c *Controller) Authenticate(password string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var req string
	if password == "" {
		req = "AUTHENTICATE"
	} else {
		req = "AUTHENTICATE " + password
	}
	_, _, err := c.makeRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) makeRequest(request string) (int, string, error) {
	id, err := c.Cmd(request)
	if err != nil {
		return 0, "", err
	}
	c.StartResponse(id)
	defer c.EndResponse(id)
	return c.ReadResponse(250)
}

func (c *Controller) getInfo(key string) (string, error) {
	_, msg, err := c.makeRequest("GETINFO " + key)
	if err != nil {
		return "", err
	}
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if parts[0] == key {
			return parts[1], nil
		}
	}
	return "", fmt.Errorf(key + " not found")
}

func (c *Controller) getInfoInt(key string) (int, error) {
	s, err := c.getInfo(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}
