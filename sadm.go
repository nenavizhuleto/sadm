package sadm

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type SAdm struct {
	Name   string
	Prefix string

	mx       sync.Mutex
	commands []command

	channel *Channel

	onConnection func(s *SAdm, c *Connection) error
}

func New(name string) *SAdm {
	return &SAdm{
		Name:   name,
		Prefix: ">>>",

		commands: make([]command, 0),

		channel: newChannel(name),

		// Default greetings message
		onConnection: func(s *SAdm, c *Connection) error {
			if err := c.Printf("welcome to %s\n", s.Name); err != nil {
				return err
			}
			if err := c.Printf("type 'help' to list available commands\n"); err != nil {
				return err
			}

			return nil
		},
	}
}

// Broadcast sends messages to main channel
func (s *SAdm) Broadcast(format string, args ...any) []error {
	return s.channel.Broadcast(format, args...)
}

func (s *SAdm) NewChannel(name string) *Channel {
	return newChannel(name)
}

// Listen starts listening on provided port
// and sets default commands: 'help' and 'exit'
func (s *SAdm) Listen(port string) error {

	s.AddCommand("help", "show help", func(c *Connection) error {
		c.Printf("available commands:\n")
		for _, cmd := range s.commands {
			c.Printf(cmd.String())
		}
		return nil
	})

	s.AddCommand("exit", "terminate session", func(c *Connection) error {
		c.Printf("goodbye\n")
		return c.Close()
	})

	ls, err := net.Listen("tcp", ":3999")
	if err != nil {
		return err
	}

	for {
		c, err := ls.Accept()
		if err != nil {
			log.Println("SAdm: ", err)
			continue
		}
		go s.handleConnection(newConnection(c))
	}
}

// AddCommand Adds command to SAdm. it panics if command already exists
func (s *SAdm) AddCommand(name string, desc string, handler Handler) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if _, ok := s.findCommand(name); ok {
		panic(fmt.Errorf("command already exists: %s", name))
	}

	cmd := newCommand(name, desc, handler)
	s.commands = append(s.commands, cmd)
}

// OnConnection callback when new connection established
func (s *SAdm) OnConnection(callback func(s *SAdm, c *Connection) error) {
	s.onConnection = callback
}

func (s *SAdm) handleConnection(c *Connection) {
	var (
		addr = c.Addr()
		err  error
	)

	s.channel.Subscribe(*c)
	defer func() {
		s.channel.Unsubscribe(addr)
		c.Close()
	}()

	if s.onConnection != nil {
		if err = s.onConnection(s, c); err != nil {
			return
		}
	}

	err = s.repl(c)
}

func (s *SAdm) findCommand(cmd string) (command, bool) {
	for _, command := range s.commands {
		if command.name == cmd {
			return command, true
		}
	}
	return command{}, false
}

func (s *SAdm) execute(cmd string, c *Connection) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if command, ok := s.findCommand(cmd); ok {
		return command.Execute(c)
	} else {
		return fmt.Errorf("unknown command: %s\n", cmd)
	}
}

func (s *SAdm) repl(c *Connection) (err error) {
	for {
		var cmd string
		if err = c.Printf("%s ", s.Prefix); err != nil {
			return
		}
		if err := c.Scanf("%s", &cmd); err != nil {
			c.Printf("%s\n", err)
			continue
		}
		if err := s.execute(cmd, c); err != nil {
			c.Printf("%s\n", err.Error())
			continue
		}
	}
}
