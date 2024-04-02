package sadm

import (
	"fmt"
)

type Handler func(*Connection) error

type Command struct {
	Name        string
	Description string
	handler     Handler
}

func NewCommand(name, desc string, handler Handler) Command {
	return Command{
		Name:        name,
		Description: desc,
		handler:     handler,
	}
}

func (c Command) Execute(connection *Connection) error {
	return c.handler(connection)
}

func (c Command) String() string {
	return fmt.Sprintf("%s\t-\t%s\n", c.Name, c.Description)
}
