package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type service struct {
	ctx     context.Context
	command []string
}

func (s *service) Run(eventChan chan<- interface{}) {
	log.Printf("Starting service %v...", s.command)

	program := s.command[0]
	argv := s.command[1:]
	cmd := exec.CommandContext(s.ctx, program, argv...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		eventChan <- err
		return
	}
	log.Printf("Service %v started.", s.command)

	eventChan <- "service started"

	err = cmd.Wait()
	if err != nil {
		log.Printf("Waiting for service %v failed: %v", s.command, err)
		eventChan <- err
		return
	}

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		log.Printf("Service %v stopped with %d.", s.command, exitCode)
		eventChan <- fmt.Errorf("process exited with %d", exitCode)
		return
	}

	log.Printf("Service %v exited normally.", s.command)

	eventChan <- "service stopped"
}
