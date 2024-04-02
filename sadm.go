package sadm

import (
	"fmt"
	"net"
	"sync"
)

type SAdm struct {
	Name   string
	Prefix string

	mx       sync.Mutex
	commands []Command

	channel *Channel

	onConnection func(s *SAdm, c *Connection) error
}

func New(name string) *SAdm {
	return &SAdm{
		Name:   name,
		Prefix: ">>>",

		commands: make([]Command, 0),

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
		return PrintHelp(c, s.commands)
	})

	s.AddCommand("exit", "terminate session", func(c *Connection) error {
		c.Printf("goodbye\n")
		return c.Close()
	})

	ls, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		c, err := ls.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(newConnection(c))
	}
}

func (s *SAdm) Use(p Plugin) {
	commands := p.Commands()
	// Usage: net ip
	// net	- plugin name
	// ip	- plugin 'ip' command
	s.AddCommand(p.Name(), p.Description(), func(c *Connection) error {
		// Read next command
		cmd, err := s.read(c)
		if err != nil {
			return err
		}

		for _, command := range commands {
			if command.Name == cmd {
				return command.Execute(c)
			}
		}

		// Handle unknown command
		return p.Help(c)
	})
}

// AddCommand Adds command to SAdm. it panics if command already exists
func (s *SAdm) AddCommand(name string, desc string, handler Handler) {

	if _, ok := s.findCommand(name); ok {
		panic(fmt.Errorf("command already exists: %s", name))
	}

	cmd := NewCommand(name, desc, handler)
	s.addCommand(cmd)
}

func (s *SAdm) addCommand(cmd Command) {
	s.mx.Lock()
	defer s.mx.Unlock()
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

	s.channel.Subscribe(c)
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

func (s *SAdm) findCommand(cmd string) (Command, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	for _, command := range s.commands {
		if command.Name == cmd {
			return command, true
		}
	}
	return Command{}, false
}

func (s *SAdm) execute(cmd string, c *Connection) error {
	if command, ok := s.findCommand(cmd); ok {
		return command.Execute(c)
	} else {
		return unknownCommandErr(cmd)
	}
}

func (s *SAdm) read(c *Connection) (string, error) {
	var cmd string
	err := c.Scanf("%s", &cmd)
	return cmd, err
}

func (s *SAdm) repl(c *Connection) (err error) {
	for {
		if err = c.Printf("%s ", s.Prefix); err != nil {
			return
		}
		cmd, err := s.read(c)
		if err != nil {
			c.Printf("%s\n", err.Error())
			continue
		}
		if err := s.execute(cmd, c); err != nil {
			c.Printf("%s\n", err.Error())
			continue
		}
	}
}
