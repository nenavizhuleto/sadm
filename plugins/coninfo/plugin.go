package coninfo

import (
	"os"

	"github.com/nenavizhuleto/sadm"
)

type coninfo struct {
	commands []sadm.Command
}

func New() sadm.Plugin {
	return &coninfo{
		commands: []sadm.Command{
			sadm.NewCommand("addr", "print client address", func(c *sadm.Connection) error {
				return c.Println(c.Addr())
			}),
			sadm.NewCommand("host", "print server address", func(c *sadm.Connection) error {
				hostname, err := os.Hostname()
				if err != nil {
					return err
				}

				return c.Println(hostname)
			}),
		},
	}
}

func (c *coninfo) Name() string {
	return "cinfo"
}

func (c *coninfo) Description() string {
	return "[c]onnection [info]rmation"
}

func (c *coninfo) Help(con *sadm.Connection) error {
	return sadm.PrintHelp(con, c.commands)
}

func (c *coninfo) Commands() []sadm.Command {
	return c.commands
}
