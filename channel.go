package sadm

import (
	"fmt"
	"sync"
)

type Channel struct {
	Name string

	mx   sync.Mutex
	cons map[string]*Connection
}

func newChannel(name string) *Channel {
	return &Channel{
		Name: name,
		cons: make(map[string]*Connection, 0),
	}
}

// Subscribe subscribes connection to the channel
func (c *Channel) Subscribe(con *Connection) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.cons[con.Addr()] = con
}

// Subscribe unsubscribes connection from the channel
func (c *Channel) Unsubscribe(addr string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	delete(c.cons, addr)
}

// Broadcast sends message to all channel subscribers
func (c *Channel) Broadcast(format string, args ...any) []error {
	c.mx.Lock()
	defer c.mx.Unlock()

	errs := make([]error, 0)
	for name, conn := range c.cons {
		if err := conn.Printf(format, args...); err != nil {
			errs = append(errs, fmt.Errorf("%s: %s", name, err))
		}
	}

	return errs
}
