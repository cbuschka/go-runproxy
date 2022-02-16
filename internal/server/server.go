package server

import (
	"context"
	"github.com/cbuschka/go-runproxy/internal/config"
	"github.com/cbuschka/go-runproxy/internal/probe"
	"github.com/cbuschka/go-runproxy/internal/proxy"
	"github.com/cbuschka/go-runproxy/internal/service"
	"log"
	"net/http"
)

func NewServer(cfg *config.Config) (*Server, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	eventChan := make(chan interface{})

	server := Server{ctx: ctx, cancelFunc: cancelFunc,
		eventChan: eventChan,
		proxy:     proxy.NewProxy(ctx, cfg.Proxy.TargetBaseUrl),
		service:   service.NewService(ctx, cfg.Service.Command),
		probe:     probe.NewProbe(ctx, cfg.Probe)}
	return &server, nil
}

type Server struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	eventChan  chan interface{}
	listenAddr string
	proxy      *proxy.Proxy
	service    *service.Service
	probe      *probe.Probe
}

func (s *Server) Run() error {

	go s.service.Run(s.eventChan)

	defer s.shutdown()

	for {
		select {
		case event := <-s.eventChan:
			log.Printf("event: %v", event)
			if err, isErr := event.(error); isErr {
				return err
			} else if "Service started" == event {
				go s.startProbe()
			} else if "Service available" == event {
				go s.startProxy()
			} else if "Service stopped" == event {
				return nil
			}
		case _ = <-s.ctx.Done():
			break
		}
	}
}

func (s *Server) shutdown() {
	log.Println("Shutting down...")
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}

func (s *Server) startProxy() {

	log.Println("Starting Proxy server on", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, s.proxy); err != nil {
		s.eventChan <- err
	}
}

func (s *Server) startProbe() {

	log.Println("Starting Probe...")
	s.probe.Watch(s.eventChan)
}
