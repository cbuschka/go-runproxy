package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsesArgs(t *testing.T) {

	cfg := NewConfig()
	err := cfg.Parse([]string{"-l", "0.0.0.0:8088",
		"-d", "localhost:8000",
		"-m", "^.*$",
		"--", "python3", "-m", "http.server"})
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "0.0.0.0:8088", cfg.ListenAddress)
	assert.Equal(t, "localhost:8000", cfg.TargetAddress)
	assert.Equal(t, "^.*$", cfg.StartupMessageMatch)
	assert.Equal(t, []string{"python3", "-m", "http.server"}, cfg.Service.Command)
}

func TestMatchPatternRequired(t *testing.T) {

	cfg := NewConfig()
	err := cfg.Parse([]string{"-l", "0.0.0.0:8088",
		"-d", "localhost:8000",
		"--", "python3", "-m", "http.server"})

	assert.Equal(t, "the required flag `-m, --match-line' was not specified", err.Error())
}

func TestServiceCommandRequired(t *testing.T) {

	cfg := NewConfig()
	err := cfg.Parse([]string{"-l", "0.0.0.0:8088",
		"-d", "localhost:8000",
		"-m", "^.*$"})

	assert.Equal(t, "the required argument `Command (at least 1 argument)` was not provided", err.Error())
}
