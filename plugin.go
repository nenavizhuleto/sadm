package sadm

type Plugin interface {
	Name() string
	Description() string
	Help(*Connection) error
	Commands() []Command
}
