package command

import (
	"github.com/cbuschka/go-runproxy/internal/config"
	"github.com/cbuschka/go-runproxy/internal/server"
)

func Run(commandLine []string) error {
	cfg, err := config.NewConfig(commandLine)
	if err != nil {
		return err
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	err = srv.Run()
	return err
}
