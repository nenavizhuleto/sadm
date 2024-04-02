package main

import (
	"github.com/nenavizhuleto/sadm"
	"log"
)

func main() {

	s := sadm.New("Example Server")

	s.AddCommand("hello", "world", func(c *sadm.Connection) error {
		return c.Printf("hello world\n")
	})

	if err := s.Listen(":3999"); err != nil {
		log.Fatalln(err)
	}
}
