package main

import (
	"context"
	"log"
	"time"

	"github.com/nenavizhuleto/sadm"
)

func main() {

	s := sadm.New("Example Server")
	ctx, cancel := context.WithCancel(context.Background())

	s.AddCommand("hello", "world", func(c *sadm.Connection) error {
		return c.Printf("hello world\n")
	})

	// Emulate context closing in some other function
	go func() {
		<-time.After(5 * time.Second)
		cancel()
	}()

	// Call s.Close() after context is done
	go func() {
		<-ctx.Done()
		log.Println("sadm close: ", s.Close())

	}()

	// Will successfully exit after 5 second timeout
	if err := s.Listen(":3999"); err != nil {
		log.Fatalln(err)
	}
}
