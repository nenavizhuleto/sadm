package sadm

import (
	"fmt"
	"net"
)

type Connection struct {
	c net.Conn
}

func newConnection(c net.Conn) *Connection {
	return &Connection{
		c: c,
	}
}

func (c *Connection) Addr() string {
	return c.c.RemoteAddr().String()
}

func (c *Connection) Printf(format string, args ...any) error {
	_, err := fmt.Fprintf(c.c, format, args...)
	return err
}
func (c *Connection) Println(args ...any) error {
	_, err := fmt.Fprintln(c.c, args...)
	return err
}

func (c *Connection) Scanf(format string, args ...any) error {
	if _, err := fmt.Fscanf(c.c, format, args...); err != nil {
		return err
	}
	return nil
}

func (c *Connection) Close() error {
	return c.c.Close()
}
