package server

import (
	"context"
	"github.com/cbuschka/go-runproxy/internal/config"
	"github.com/cbuschka/go-runproxy/internal/console"
	"github.com/cbuschka/go-runproxy/internal/proxy"
	"github.com/cbuschka/go-runproxy/internal/service"
	"os"
	"os/signal"
)

func NewServer(cfg *config.Config) (*Server, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	eventChan := make(chan interface{})

	svc, err := service.NewService(ctx, cfg.Service.Command, cfg.StartupMessageMatch)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	prx, err := proxy.NewProxy(ctx, cfg)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	server := Server{ctx: ctx, cancelFunc: cancelFunc,
		eventChan: eventChan,
		proxy:     prx,
		service:   svc}
	return &server, nil
}

type Server struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	eventChan  chan interface{}
	proxy      *proxy.TcpProxyStrategy
	service    *service.Service
}

func (s *Server) Run() error {

	s.service.Start(s.eventChan)

	s.installSigHandler()

	defer s.shutdown()

	for {
		select {
		case event := <-s.eventChan:
			console.Debugf("Event \"%v\" seen.", event)
			if err, isErr := event.(error); isErr {
				return err
			} else if "service started" == event {
				console.Info("Service started. Waiting to come up.")
			} else if "service available" == event || "startup message seen" == event {
				console.Info("Service available. Starting proxy.")
				go s.startProxy()
			} else if "proxy started" == event {
				break
			} else if "service stopped" == event {
				return nil
			} else if "shutdown" == event {
				break
			} else {
				console.Debugf("Event %s unknown.", event)
			}
		case _ = <-s.ctx.Done():
			break
		}
	}
}

func (s *Server) installSigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		for _ = range c {
			s.shutdown()
		}
	}()
}

func (s *Server) shutdown() {
	console.Info("Shutting down...")

	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	if s.service != nil {
		console.Info("Killing service...")
		s.service.Kill()
	}
}

func (s *Server) startProxy() {
	s.proxy.Start(s.ctx, s.eventChan)
}
