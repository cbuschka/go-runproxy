package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadsCompleteConfigNoOverrides(t *testing.T) {

	cmdline := []string{"", "-c", "../../example-config.yml"}

	cfg, err := NewConfig(cmdline)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "runproxy/1", cfg.Version)
	assert.Equal(t, "0.0.0.0:8088", cfg.Proxy.Http.ListenAddress)
	assert.Equal(t, "http://localhost:8000", cfg.Proxy.Http.TargetBaseUrl)
	assert.Equal(t, []string{"python3", "-m", "http.server"}, cfg.Service.Command)
	assert.Equal(t, "^.*$", cfg.Service.StartupMessageMatch)
	assert.Equal(t, "http://localhost:8000", cfg.Healthcheck.EndpointAddress)
	assert.Equal(t, []string{"curl", "-sfL", "-o", "/dev/null", "localhost:8000"}, cfg.Healthcheck.Command)
	assert.Equal(t, 500, cfg.Healthcheck.CheckIntervalMillis)
	assert.Equal(t, 30000, cfg.Healthcheck.RecheckIntervalMillis)
}
