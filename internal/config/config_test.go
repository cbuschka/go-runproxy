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
		"-t", "30000",
		"--", "python3", "-m", "http.server"})
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "0.0.0.0:8088", cfg.ListenAddress)
	assert.Equal(t, "localhost:8000", cfg.TargetAddress)
	assert.Equal(t, "^.*$", cfg.StartupMessageMatch)
	assert.Equal(t, uint(30000), cfg.StartupTimeoutMillis)
	assert.Equal(t, []string{"python3", "-m", "http.server"}, cfg.Service.Command)
}
