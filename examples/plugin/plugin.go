package main

import (
	"log"

	"github.com/nenavizhuleto/sadm"
	"github.com/nenavizhuleto/sadm/plugins/coninfo"
)

func main() {
	s := sadm.New("Plugins")

	s.Use(coninfo.New())

	if err := s.Listen(":3999"); err != nil {
		log.Fatal(err)
	}
}
