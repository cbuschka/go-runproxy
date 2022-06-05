package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"syscall"
	"testing"
)

const fullConfig = `version: "runproxy/1"
proxy:
  http:
    listenAddress: 0.0.0.0:8088
    targetBaseUrl: http://localhost:8000
service:
  command: [ "python3", "-m", "http.server" ]
  startupMessageMatch: '^.*$'
healthcheck:
  command: [ "curl", "-sfL", "-o", "/dev/null", "localhost:8000" ]
  http:
    endpointAddress: "http://localhost:8000"
  checkIntervalMillis: 500
  recheckIntervalMillis: 30000
`

func TestLoadsCompleteConfigNoOverrides(t *testing.T) {

	tempConfigFile, err := ioutil.TempFile("", "config.yml")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer syscall.Unlink(tempConfigFile.Name())
	err = ioutil.WriteFile(tempConfigFile.Name(), []byte(fullConfig), 0644)
	if err != nil {
		t.Fatal(err)
		return
	}

	cmdline := []string{"", "-c", tempConfigFile.Name()}

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
	assert.Equal(t, "http://localhost:8000", cfg.Healthcheck.Http.EndpointAddress)
	assert.Equal(t, []string{"curl", "-sfL", "-o", "/dev/null", "localhost:8000"}, cfg.Healthcheck.Command)
	assert.Equal(t, 500, cfg.Healthcheck.CheckIntervalMillis)
	assert.Equal(t, 30000, cfg.Healthcheck.RecheckIntervalMillis)
}
