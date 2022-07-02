package command

import (
	"github.com/cbuschka/go-runproxy/internal/build"
	"github.com/cbuschka/go-runproxy/internal/config"
	"github.com/cbuschka/go-runproxy/internal/console"
	"github.com/cbuschka/go-runproxy/internal/server"
	"os"
)

func Run(commandLine []string) error {
	cfg := config.NewConfig()
	err := cfg.Parse(commandLine[1:])
	if err != nil {
		cfg.PrintUsage(os.Stderr)
		return err
	}

	printRunInfo()

	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	err = srv.Run()
	return err
}

func printRunInfo() {
	console.Infof("go-runproxy %s (https://github.com/cbuschka/go-runproxy)", build.Version)
}
