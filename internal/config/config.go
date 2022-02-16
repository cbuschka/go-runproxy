package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type HealthcheckConfig struct {
	Command               []string
	EndpointAddress       string
	CheckIntervalMillis   time.Duration
	RecheckIntervalMillis time.Duration
}

type ProxyConfig struct {
	ListenAddress string `yaml:"listenAddress"`
	TargetBaseUrl string `yaml:"targetBaseUrl"`
}

type ServiceConfig struct {
	Command []string `yaml:"command"`
}

type Config struct {
	Version     string            `yaml:"version"`
	Proxy       ProxyConfig       `yaml:"proxy"`
	Service     ServiceConfig     `yaml:"service"`
	Healthcheck HealthcheckConfig `yaml:"healthcheck"`
}

func (c *Config) overrideFromCommandLine(commandLine []string, cmdlineCfg *cmdlineConfig) error {

	serviceCommand := extractServiceCommand(commandLine)
	if len(serviceCommand) > 0 {
		c.Service.Command = serviceCommand
	}

	if cmdlineCfg.ListenAddress != "" {
		log.Printf("Overriding listen address from command line: %s", cmdlineCfg.ListenAddress)
		c.Proxy.ListenAddress = cmdlineCfg.ListenAddress
	}

	if cmdlineCfg.ServiceCommand != "" {
		log.Printf("Overriding service command from command line: %s", cmdlineCfg.ServiceCommand)
		c.Service.Command = strings.Split(cmdlineCfg.ServiceCommand, " ")
	}

	return nil
}

func NewConfig(commandLine []string) (*Config, error) {
	cmdlnCfg, err := parseCommandline(commandLine)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Version: "1",
		Proxy:   ProxyConfig{ListenAddress: ":8080"},
		Service: ServiceConfig{Command: []string{"python3", "-m", "http.server"}},
		Healthcheck: HealthcheckConfig{Command: []string{"curl", "-sLf", "http://localhost:8000"},
			CheckIntervalMillis:   300,
			RecheckIntervalMillis: 30000,
		},
	}

	if cmdlnCfg.ConfigFile != "" {
		log.Printf("Loading config from %s...", cmdlnCfg.ConfigFile)

		cfgBytes, err := os.ReadFile(cmdlnCfg.ConfigFile)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(cfgBytes, &cfg)
		if err != nil {
			return nil, err
		}

		if cfg.Version != "runproxy/1" {
			return nil, fmt.Errorf("invalid version, expected runproxy/1")
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
