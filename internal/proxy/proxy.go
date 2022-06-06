package proxy

import (
	"context"
	"fmt"
	"github.com/cbuschka/go-runproxy/internal/config"
	"log"
)

type ProxyStrategy interface {
	Start(ctx context.Context, eventChan chan interface{})
}

func NewProxy(ctx context.Context, cfg config.ProxyConfig) (*Proxy, error) {
	var strategy ProxyStrategy
	if cfg.Http != nil {
		strategy = ProxyStrategy(&HttpProxyStrategy{listenAddress: cfg.Http.ListenAddress, targetBaseUrl: cfg.Http.TargetBaseUrl})
	} else if cfg.Tcp != nil {
		strategy = ProxyStrategy(&TcpProxyStrategy{listenAddress: cfg.Tcp.ListenAddress, targetEndpointAddress: cfg.Tcp.TargetEndpointAddress})
	} else {
		return nil, fmt.Errorf("neither tcp nor http proxy configured")
	}

	return &Proxy{ctx: ctx, strategy: strategy}, nil
}

type Proxy struct {
	ctx      context.Context
	strategy ProxyStrategy
}

func (p *Proxy) Start(eventChan chan interface{}) {

	log.Printf("Starting proxy %s...", p)
	p.strategy.Start(p.ctx, eventChan)
}

func (p *Proxy) Kill() {

}

func (p *Proxy) String() string {
	return fmt.Sprintf("Proxy{strategy=%v}", p.strategy)
}
