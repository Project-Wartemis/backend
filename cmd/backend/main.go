package main

import (
	"flag"
	"log"

	"github.com/Project-Wartemis/pw-backend/pkg/backend"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.Println("Execute main")

	communication.RunBackend(addr)
}
