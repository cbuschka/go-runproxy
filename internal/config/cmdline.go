package config

import (
	optsPkg "github.com/jpillora/opts"
)

type cmdlineConfig struct {
	ConfigFile     string `opts:"help=Config file in yaml format,short=c"`
	ListenAddress  string `opts:"help=Proxy listen address,short=l"`
	ServiceCommand string `opts:"help=Command to start service,short=e"`
}

func parseCommandline(args []string) (*cmdlineConfig, error) {
	cfg := cmdlineConfig{}
	opts := optsPkg.New(&cfg)
	_, err := opts.ParseArgsError(args)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
