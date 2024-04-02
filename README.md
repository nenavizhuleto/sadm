# Simple [S]ocket [Adm]inistration (SAdm) Server

## Usage 

### Simple commands

```go
s := sadm.New("Example Server")

s.AddCommand("hello", "world", func(c *sadm.Connection) error {
    return c.Printf("hello world\n")
})

if err := s.Listen(":3999"); err != nil {
    log.Fatalln(err)
}
```

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

See [examples/server](examples/server/server.go) for demo

### Broadcasting and Channels

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

>>> room
hello from room!
hello from main!
hello from room!
hello from main!
hello from room!
...
```

See [examples/broadcast](examples/broadcast/broadcast.go) for demo

## Installation

```sh
go get github.com/nenavizhuleto/sadm
```
