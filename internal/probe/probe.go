package probe

import (
	"context"
	"fmt"
	"github.com/cbuschka/go-runproxy/internal/config"
	"log"
	"os"
	"os/exec"
	"time"
)

type Probe struct {
	command        []string
	ctx            context.Context
	checkTimeout   time.Duration
	recheckTimeout time.Duration
}

func NewProbe(ctx context.Context, cfg config.ProbeConfig) *Probe {
	prb := Probe{ctx: ctx,
		command:        cfg.Command,
		checkTimeout:   cfg.CheckIntervalMillis,
		recheckTimeout: cfg.RecheckIntervalMillis}
	return &prb
}

func (p *Probe) Watch(eventChan chan<- interface{}) {

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

func (p *Probe) isAvailable() (bool, error) {

	log.Println("Checking if service is available...")

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
