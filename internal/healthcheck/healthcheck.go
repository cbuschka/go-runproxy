package healthcheck

import (
	"context"
	"fmt"
	"github.com/cbuschka/go-runproxy/internal/config"
	"log"
	"time"
)

type Healthcheck struct {
	ctx            context.Context
	checkTimeout   time.Duration
	recheckTimeout time.Duration
	strategy       HealthcheckStrategy
}

func NewHealthcheck(ctx context.Context, cfg config.HealthcheckConfig) *Healthcheck {
	var strategy HealthcheckStrategy
	if cfg.Command != nil && len(cfg.Command) > 0 {
		strategy = HealthcheckStrategy(&CommandHealthcheckStrategy{command: cfg.Command})
	} else if cfg.Http != nil {
		strategy = HealthcheckStrategy(&HttpHealthcheckStrategy{endpointAddress: cfg.Http.EndpointAddress})
	} else {
		strategy = HealthcheckStrategy(&DummyHealthcheckStrategy{})
	}

	log.Printf("Using healthcheck strategy %v.", strategy)

	prb := Healthcheck{ctx: ctx,
		strategy:       strategy,
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

	return p.strategy.IsAvailable(p.ctx)
}
