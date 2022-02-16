package config

import (
	"time"
)

type ProbeConfig struct {
	Command               []string
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
	Probe   ProbeConfig
}

func NewConfig(commandLine []string) (*Config, error) {
	cfg := Config{
		Proxy:   ProxyConfig{ListenAddress: ":8080"},
		Service: ServiceConfig{Command: []string{"python3", "-m", "http.server"}},
		Probe: ProbeConfig{Command: []string{"curl", "-sLf", "http://localhost:8000"},
			CheckIntervalMillis:   300 * time.Millisecond,
			RecheckIntervalMillis: 30 * time.Second,
		},
	}

	cfg.Service.Command = extractServiceCommand(commandLine)

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
