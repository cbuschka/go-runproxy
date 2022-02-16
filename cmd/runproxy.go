package main

import (
	"github.com/cbuschka/go-runproxy/internal/command"
	"log"
	"os"
)

func main() {
	err := command.Run(os.Args)
	if err != nil {
		log.Printf("failed: %v", err)
		os.Exit(1)
	}
}
