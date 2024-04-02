package perfmon

import "github.com/nenavizhuleto/sadm"

type perfmon struct {
	commands []sadm.Command
}

func New() sadm.Plugin {
	return &perfmon{
		commands: []sadm.Command{
			sadm.NewCommand("cpu", "print cpu info", cpuInfo),
			sadm.NewCommand("stat", "print cpu stats", stat),
		},
	}
}

func (p *perfmon) Name() string {
	return "perfmon"
}

func (p *perfmon) Description() string {
	return "[perf]ormance [mon]itoring"
}

func (p *perfmon) Help(c *sadm.Connection) error {
	return sadm.PrintHelp(c, p.commands)
}

func (p *perfmon) Commands() []sadm.Command {
	return p.commands
}
