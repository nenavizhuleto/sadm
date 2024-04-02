package main

import (
	"log"

	"github.com/nenavizhuleto/sadm"
	"github.com/nenavizhuleto/sadm/plugins/perfmon"
)

func main() {
	s := sadm.New("Performance monitor")

	s.Use(perfmon.New())

	if err := s.Listen(":3999"); err != nil {
		log.Fatal(err)
	}
}
