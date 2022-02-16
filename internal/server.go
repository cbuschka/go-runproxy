package internal

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	eventChan  chan interface{}
	listenAddr string
	proxy      *proxy
	service    *service
	probe      *probe
}

func Run() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	eventChan := make(chan interface{})
	proxy := proxy{targetBaseUrl: "http://localhost:8000"}
	service := service{ctx: ctx,
		command: []string{"python3", "-m", "http.server"}}
	probe := probe{ctx: ctx,
		command:      []string{"curl", "-sLf", "-o", "/dev/null", "http://localhost:8000"},
		checkTimeout: 300 * time.Millisecond, recheckTimeot: 30 * time.Second}
	server := Server{ctx: ctx, cancelFunc: cancelFunc,
		eventChan: eventChan,
		proxy:     &proxy,
		service:   &service,
		probe:     &probe}
	if err := server.init(); err != nil {
		return err
	}
	return server.Run()
}

func (s *Server) init() error {
	var addr = flag.String("addr", "127.0.0.1:8080", "The addr of the application.")
	flag.Parse()

	s.listenAddr = *addr

	return nil
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
