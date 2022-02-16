package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type HealthcheckConfig struct {
	Command               []string
	EndpointAddress       string
	CheckIntervalMillis   time.Duration
	RecheckIntervalMillis time.Duration
}

type ProxyConfig struct {
	ListenAddress string
	TargetBaseUrl string
}

type ServiceConfig struct {
	Command []string
}

type Config struct {
	Proxy   ProxyConfig
	Service ServiceConfig
	Probe   HealthcheckConfig
}

func (c Config) overrideFromCommandLine(commandLine []string, cfg *cmdlineConfig) error {

	serviceCommand := extractServiceCommand(commandLine)
	if len(serviceCommand) > 0 {
		c.Service.Command = serviceCommand
	}

	if cfg.ListenAddress != "" {
		c.Proxy.ListenAddress = cfg.ListenAddress
	}

	return nil
}

func NewConfig(commandLine []string) (*Config, error) {
	cmdlnCfg, err := parseCommandline(commandLine)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Proxy:   ProxyConfig{ListenAddress: ":8080"},
		Service: ServiceConfig{Command: []string{"python3", "-m", "http.server"}},
		Probe: HealthcheckConfig{Command: []string{"curl", "-sLf", "http://localhost:8000"},
			CheckIntervalMillis:   300 * time.Millisecond,
			RecheckIntervalMillis: 30 * time.Second,
		},
	}

	if cmdlnCfg.ConfigFile != "" {
		cfgBytes, err := os.ReadFile(cmdlnCfg.ConfigFile)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(cfgBytes, &cfg)
		if err != nil {
			return nil, err
		}
	}

	err = cfg.overrideFromCommandLine(commandLine, cmdlnCfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func extractServiceCommand(commandLine []string) []string {
	cmd := []string{}
	doubleDashSeen := false
	for _, arg := range commandLine {
		if doubleDashSeen {
			cmd = append(cmd, arg)
		} else if arg == "--" {
			doubleDashSeen = true
		}
	}

	return cmd
}
