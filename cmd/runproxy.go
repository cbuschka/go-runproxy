package main

import (
	"fmt"
	"github.com/cbuschka/go-runproxy/internal/command"
	"os"
)

func main() {
	err := command.Run(os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed: %v", err)
		os.Exit(1)
	}
}
