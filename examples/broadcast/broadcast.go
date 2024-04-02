package main

import (
	"log"
	"time"

	"github.com/nenavizhuleto/sadm"
)

func main() {
	s := sadm.New("Broadcast")

	ch := s.NewChannel("room")

	s.AddCommand("room", "another broadcast", func(c *sadm.Connection) error {
		ch.Subscribe(c)

		go func() {
			<-time.After(5 * time.Second)
			ch.Unsubscribe(c.Addr())
		}()

		return nil
	})

	go func() {
		if err := s.Listen(":3999"); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		// Broadcast to main channel
		s.Broadcast("hello from main!\n")

		// Broadcast to specific channel
		ch.Broadcast("hello from room!\n")

		time.Sleep(2 * time.Second)
	}
}
