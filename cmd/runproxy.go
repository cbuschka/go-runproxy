package main

import (
	"github.com/cbuschka/go-runproxy/internal"
	"log"
	"os"
)

func main() {
	err := internal.Run()
	if err != nil {
		log.Printf("failed: %v", err)
		os.Exit(1)
	}
}
