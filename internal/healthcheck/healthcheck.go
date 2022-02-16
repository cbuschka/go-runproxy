package healthcheck

import (
	"context"
	"fmt"
	"github.com/cbuschka/go-runproxy/internal/config"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Healthcheck struct {
	command         []string
	endpointAddress string
	ctx             context.Context
	checkTimeout    time.Duration
	recheckTimeout  time.Duration
}

func NewHealthcheck(ctx context.Context, cfg config.HealthcheckConfig) *Healthcheck {
	prb := Healthcheck{ctx: ctx,
		command:         cfg.Command,
		endpointAddress: cfg.EndpointAddress,
		checkTimeout:   time.Duration(cfg.CheckIntervalMillis) * time.Millisecond,
		recheckTimeout: time.Duration(cfg.RecheckIntervalMillis) * time.Millisecond}
	return &prb
}

func (p *Healthcheck) Watch(eventChan chan<- interface{}) {

	serviceAvailable := false
	checkTimeout := p.checkTimeout
	for {
		select {
		case <-time.After(checkTimeout):
			available, err := p.isAvailable()
			if err != nil {
				eventChan <- err
			} else if available {
				checkTimeout = p.recheckTimeout
				if !serviceAvailable {
					serviceAvailable = true
					eventChan <- "service available"
				}
			} else {
				checkTimeout = p.checkTimeout
				if serviceAvailable {
					serviceAvailable = false
					eventChan <- fmt.Errorf("service not available")
				}
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Healthcheck) isAvailable() (bool, error) {

	log.Println("Checking if service is available...")

	if p.command != nil && len(p.command) > 0 {
		return p.isAvailableByCommand()
	} else if strings.HasPrefix(p.endpointAddress, "http://") {
		return p.isAvailableViaHttp()
	} else {
		return false, fmt.Errorf("unsupported endpoint %s", p.endpointAddress)
	}
}

func (p *Healthcheck) isAvailableViaHttp() (bool, error) {
	endpointUrl, err := url.Parse(p.endpointAddress)
	if err != nil {
		return false, err
	}

	client := &http.Client{}

	req := http.Request{
		URL:    endpointUrl,
		Method: "GET",
	}

	resp, err := client.Do(&req)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}

func (p *Healthcheck) isAvailableByCommand() (bool, error) {
	cmd := exec.CommandContext(p.ctx, p.command[0], p.command[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return false, err
	}

	err = cmd.Wait()
	if err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		if exitCode != 0 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
