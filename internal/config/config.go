package config

import (
	"github.com/jessevdk/go-flags"
	"io"
)

type Config struct {
	Verbose             bool   `short:"v" long:"verbose" description:"enable verbose output"`
	ListenAddress       string `required:"yes" short:"l" long:"listen-address" description:"listen address ip:port"`
	TargetAddress       string `required:"yes" short:"d" long:"destination-address" description:"destination address ip:port"`
	StartupMessageMatch string `required:"yes" short:"m" long:"match-line" description:"regex for matching startup message line"`
	Service             struct {
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
