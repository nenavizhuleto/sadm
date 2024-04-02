# Simple [S]ocket [Adm]inistration (SAdm) Server

## Usage 

### Simple commands

You can add and execute simple commands (aka. functions)

Example:
```sh
$ nc localhost 3999

welcome to Example Server
type 'help' to list available commands
>>> help
available commands:
hello	-	world
help	-	show help
exit	-	terminate session
>>> exit
goodbye
```

```go
s := sadm.New("Example Server")

s.AddCommand("hello", "world", func(c *sadm.Connection) error {
    return c.Printf("hello world\n")
})

if err := s.Listen(":3999"); err != nil {
    log.Fatalln(err)
}
```

See [examples/server](examples/server/server.go) for demo

### Broadcasting and Channels

You can broadcast messages to all clients connected via `SAdm.Broadcast()`

Or you can create separate `Channel` and subscribe to it via `Command`

Example:
```sh
$ nc localhost 3999

welcome to Broadcast
type 'help' to list available commands
hello from main!
hello from main!
hello from main!
hello from main!
hello from main!
...

# subscribing to channel
>>> room
hello from room!
hello from main!
hello from room!
hello from main!
hello from room!
...
```

```go
s := sadm.New("Broadcast")

// Broadcast to main channel
s.Broadcast("hello from main!\n")

// Create new channel
ch := s.NewChannel("room")

s.AddCommand("room", "another broadcast", func(c *sadm.Connection) error {
	// Subscribe to channel via Command
	ch.Subscribe(*c)

	
	go func() {
		<-time.After(5 * time.Second)
		// Unsubscribe after 5 seconds
		ch.Unsubscribe(c.Addr())
	}()

	return nil
})

// Broadcast to specific channel
ch.Broadcast("hello from room!\n")
```

See [examples/broadcast](examples/broadcast/broadcast.go) for demo

### Plugins

Plugins provide already implemented commands.

You can run plugin's commands via `plugin_name plugin_command`

Example:
```sh
$ nc localhost 3999

welcome to Plugins
type 'help' to list available commands
>>> help
available commands:
	cinfo	-	[c]onnection [info]rmation
	help	-	show help
	exit	-	terminate session
>>> cinfo help
available commands:
	addr	-	print client address
	host	-	print server address
>>> cinfo addr
127.0.0.1:44766
>>> exit
goodbye
```

```go

// Import plugin
import "github.com/nenavizhuleto/sadm/plugins/coninfo"

...

s := sadm.New("Plugins")

// Use plugin
s.Use(coninfo.New())

if err := s.Listen(":3999"); err != nil {
	log.Fatal(err)
}
```

See [examples/plugin](examples/plugin/plugin.go) for demo

## Installation

```sh
go get github.com/nenavizhuleto/sadm
```
