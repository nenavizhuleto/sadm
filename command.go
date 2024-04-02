package sadm

import (
	"fmt"
)

type Handler func(*Connection) error

type command struct {
	name        string
	description string
	handler     Handler
}

func newCommand(name, desc string, handler Handler) command {
	return command{
		name:        name,
		description: desc,
		handler:     handler,
	}
}

func (c command) Execute(connection *Connection) error {
	return c.handler(connection)
}

func (c command) String() string {
	return fmt.Sprintf("%s\t-\t%s\n", c.name, c.description)
}
