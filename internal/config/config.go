package config

import (
	"github.com/jessevdk/go-flags"
	"io"
)

type Config struct {
	ListenAddress        string `short:"l" long:"listen-address" description:"listen address ip:port"`
	TargetAddress        string `short:"d" long:"destination-address" description:"destination address ip:port"`
	StartupMessageMatch  string `short:"m" long:"match-line" description:"regex for matching startup message line"`
	StartupTimeoutMillis uint   `short:"t" long:"startup-timeout" description:"max time to wait for startup finished in millis"`
	Service              struct {
		Command []string `required:"yes" description:"the downstream service to start"`
	} `positional-args:"yes" `
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) PrintUsage(out io.Writer) {
	parser := flags.NewParser(c, flags.Default)
	parser.WriteHelp(out)
}

func (c *Config) Parse(commandLine []string) error {

	parser := flags.NewParser(c, flags.PassDoubleDash)
	_, err := parser.ParseArgs(commandLine)
	if err != nil {
		return err
	}

	return nil
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
