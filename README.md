# Simple [S]ocket [Adm]inistration (SAdm) Server

## Usage 

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

## Installation

```sh
go get github.com/nenavizhuleto/sadm
```
