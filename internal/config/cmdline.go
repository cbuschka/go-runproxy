package config

import (
	optsPkg "github.com/jpillora/opts"
)

type cmdlineConfig struct {
	ConfigFile    string `opts:"help=Config file in yaml format" opts:"short=c"`
	ListenAddress string `opts:"help=Proxy listen address"`
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
