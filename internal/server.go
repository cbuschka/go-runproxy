package internal

import (
	"flag"
	"log"
	"net/http"
)

type Server struct{}

func Run() error {
	server := Server{}
	return server.Run()
}

func (s *Server) Run() error {
	var addr = flag.String("addr", "127.0.0.1:8080", "The addr of the application.")
	flag.Parse()

	handler := &proxy{}

	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		return err
	}

	return nil
}
