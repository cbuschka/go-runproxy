package command

import (
	"github.com/cbuschka/go-runproxy/internal/build"
	"github.com/cbuschka/go-runproxy/internal/config"
	"github.com/cbuschka/go-runproxy/internal/server"
	"log"
	"os"
)

func Run(commandLine []string) error {
	log.SetOutput(os.Stderr)

	cfg := config.NewConfig()
	err := cfg.Parse(commandLine[1:])
	if err != nil {
		cfg.PrintUsage(os.Stderr)
		return err
	}

	log.Printf("go-runproxy (https://github.com/cbuschka/go-runproxy)")
	log.Printf("Build version: %s", build.Version)
	log.Printf("Build timestamp: %s", build.Timestamp)
	log.Printf("Build commitish: %s", build.Commitish)
	log.Printf("Build os/arch: %s/%s", build.Os, build.Arch)

	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	err = srv.Run()
	return err
}
