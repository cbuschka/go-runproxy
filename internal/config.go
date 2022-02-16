package internal

import (
	"time"
)

type ProbeConfig struct {
	command               []string
	checkIntervalMillis   time.Duration
	recheckIntervalMillis time.Duration
}

type ProxyConfig struct {
	listenAddress string
	targetBaseUrl string
}

type ServiceConfig struct {
	command []string
}

type Config struct {
	proxy   ProxyConfig
	service ServiceConfig
	probe   ProbeConfig
}

func NewConfig(commandLine []string) (*Config, error) {
	config := Config{
		proxy:   ProxyConfig{listenAddress: ":8080"},
		service: ServiceConfig{command: []string{"python3", "-m", "http.server"}},
		probe: ProbeConfig{command: []string{"curl", "-sLf", "http://localhost:8000"},
			checkIntervalMillis:   300 * time.Millisecond,
			recheckIntervalMillis: 30 * time.Second,
		},
	}

	config.service.command = extractServiceCommand(commandLine)

	return &config, nil
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
