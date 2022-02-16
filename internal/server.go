package internal

import (
	"context"
	"log"
	"net/http"
)

func NewServer(config *Config) (*Server, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	eventChan := make(chan interface{})
	proxy := proxy{targetBaseUrl: config.proxy.targetBaseUrl}
	service := service{ctx: ctx, command: config.service.command}
	probe := probe{ctx: ctx,
		command:        config.probe.command,
		checkTimeout:   config.probe.checkIntervalMillis,
		recheckTimeout: config.probe.recheckIntervalMillis}

	server := Server{ctx: ctx, cancelFunc: cancelFunc,
		eventChan: eventChan,
		proxy:     &proxy,
		service:   &service,
		probe:     &probe}
	return &server, nil
}

type Server struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	eventChan  chan interface{}
	listenAddr string
	proxy      *proxy
	service    *service
	probe      *probe
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
			} else if "service started" == event {
				go s.startProbe()
			} else if "service available" == event {
				go s.startProxy()
			} else if "service stopped" == event {
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

	log.Println("Starting proxy server on", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, s.proxy); err != nil {
		s.eventChan <- err
	}
}

func (s *Server) startProbe() {

	log.Println("Starting probe...")
	s.probe.Watch(s.eventChan)
}
