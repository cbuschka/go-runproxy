package server

import (
	"context"
	"github.com/cbuschka/go-runproxy/internal/config"
	"github.com/cbuschka/go-runproxy/internal/healthcheck"
	"github.com/cbuschka/go-runproxy/internal/proxy"
	"github.com/cbuschka/go-runproxy/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func NewServer(cfg *config.Config) (*Server, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	eventChan := make(chan interface{})

	server := Server{ctx: ctx, cancelFunc: cancelFunc,
		eventChan:  eventChan,
		listenAddr: cfg.Proxy.ListenAddress,
		proxy:      proxy.NewProxy(ctx, cfg.Proxy.TargetBaseUrl),
		service:    service.NewService(ctx, cfg.Service.Command),
		probe:      healthcheck.NewHealthcheck(ctx, cfg.Healthcheck)}
	return &server, nil
}

type Server struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	eventChan  chan interface{}
	listenAddr string
	proxy      *proxy.Proxy
	service    *service.Service
	probe      *healthcheck.Healthcheck
}

func (s *Server) Run() error {

	s.service.Start(s.eventChan)

	s.installSigHandler()

	defer s.shutdown()

	for {
		select {
		case event := <-s.eventChan:
			log.Printf("event: %v", event)
			if err, isErr := event.(error); isErr {
				return err
			} else if "service started" == event {
				go s.startProbe()
			} else if "service available" == event {
				go s.startProxy()
			} else if "service stopped" == event {
				return nil
			} else if "shutdown" == event {
				break
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
	log.Println("Shutting down...")

	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	if s.service != nil {
		log.Println("Killing service...")
		s.service.Kill()
	}
}

func (s *Server) startProxy() {

	log.Println("Starting proxy server on", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, s.proxy); err != nil {
		s.eventChan <- err
	}
}

func (s *Server) startProbe() {

	log.Println("Starting healthcheck...")
	s.probe.Watch(s.eventChan)
}
