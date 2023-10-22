package netxclient

import (
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	DefaultNetwork      = "tcp"
	DefaultDialTimeout  = 6 * time.Second
	DefaultWriteTimeout = 6 * time.Second
	DefaultReadTimeout  = 6 * time.Second
	DefaultRetryDelay   = 2 * time.Second
	DefaultReadSize     = 1024 * 2
	DefaultMaxRetries   = 3
)

type Client struct {
	conn net.Conn
	conf Config
}

func NewClient(address string, conf Config) (*Client, error) {
	var (
		err  error
		conn net.Conn
	)

	if conf.DialTimeout == 0 {
		conf.DialTimeout = DefaultDialTimeout
	}

	if conf.RetryDelay == 0 {
		conf.RetryDelay = DefaultRetryDelay
	}

	if conf.MaxRetries == 0 {
		conf.MaxRetries = DefaultMaxRetries
	}

	if len(conf.Network) == 0 {
		conf.Network = DefaultNetwork
	}

	for i := 0; i < conf.MaxRetries; i++ {
		conn, err = net.DialTimeout(conf.Network, address, conf.DialTimeout)
		if err == nil {
			break
		}

		time.Sleep(conf.RetryDelay)
	}

	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, conf: conf}, nil
}

func (c *Client) Send(data []byte) error {
	if c.conn == nil {
		return errors.New("connection is not established")
	}

	c.conn.SetReadDeadline(time.Now().Add(c.writeTimeout()))
	_, err := c.conn.Write(data)
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			err = c.retryWrite(data)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (c *Client) Receive() ([]byte, error) {
	if c.conn == nil {
		return nil, errors.New("connection is not established")
	}

	c.conn.SetReadDeadline(time.Now().Add(c.readTimeout()))
	buf := make([]byte, c.readSize())
	n, err := c.conn.Read(buf)
	if err != nil {
		err = c.retryRead(buf)
		if err != nil {
			return nil, err
		}
	}
	return buf[:n], nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) retryWrite(data []byte) error {
	for i := 0; i < c.maxRetries(); i++ {
		time.Sleep(c.retryTimeout())
		conn, err := net.DialTimeout(c.network(), c.conn.RemoteAddr().String(), c.dialTimeout())
		if err == nil {
			c.conn = conn
			c.conn.SetReadDeadline(time.Now().Add(c.writeTimeout()))
			_, err = c.conn.Write(data)
			if err == nil {
				return nil
			}
		}
	}
	return fmt.Errorf("%s failed to send data after %d retries: %w", c.conn.RemoteAddr().String(), c.maxRetries(), errors.New("connection is closed"))
}

func (c *Client) retryRead(buf []byte) error {
	for i := 0; i < c.maxRetries(); i++ {
		time.Sleep(c.retryTimeout())
		conn, err := net.DialTimeout(c.network(), c.conn.RemoteAddr().String(), c.dialTimeout())
		if err == nil {
			c.conn = conn
			c.conn.SetReadDeadline(time.Now().Add(c.readTimeout()))
			n, err := c.conn.Read(buf)
			if err == nil {
				return nil
			}
			if n > 0 {
				return nil
			}
		}
	}
	return fmt.Errorf("%s failed to receive data after %d retries: %w", c.conn.RemoteAddr().String(), c.maxRetries(), errors.New("connection is closed"))
}

func (c *Client) network() string {
	if len(c.conf.Network) > 0 {
		return c.conf.Network
	}
	return DefaultNetwork
}

func (c *Client) dialTimeout() time.Duration {
	if c.conf.DialTimeout != 0 {
		return c.conf.DialTimeout
	}
	return DefaultDialTimeout
}

func (c *Client) writeTimeout() time.Duration {
	if c.conf.WriteTimeout != 0 {
		return c.conf.WriteTimeout
	}
	return DefaultWriteTimeout
}

func (c *Client) readTimeout() time.Duration {
	if c.conf.ReadTimeout != 0 {
		return c.conf.ReadTimeout
	}
	return DefaultReadTimeout
}

func (c *Client) retryTimeout() time.Duration {
	if c.conf.RetryDelay != 0 {
		return c.conf.RetryDelay
	}
	return DefaultRetryDelay
}

func (c *Client) maxRetries() int {
	if c.conf.MaxRetries != 0 {
		return c.conf.MaxRetries
	}
	return DefaultMaxRetries
}

func (c *Client) readSize() int {
	if c.conf.ReadSize != 0 {
		return c.conf.ReadSize
	}
	return DefaultReadSize
}
