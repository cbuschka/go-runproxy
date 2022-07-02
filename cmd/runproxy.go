package main

import (
	"github.com/cbuschka/go-runproxy/internal/command"
	"github.com/cbuschka/go-runproxy/internal/console"
	"os"
)

func main() {
	err := command.Run(os.Args)
	if err != nil {
		console.Infof("failed: %v", err)
		os.Exit(1)
	}
}
