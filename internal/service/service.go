package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Service struct {
	ctx     context.Context
	command []string

	cmd *exec.Cmd
}

func NewService(ctx context.Context, command []string) *Service {
	return &Service{ctx: ctx, command: command}
}

func (s *Service) run(eventChan chan<- interface{}) {
	log.Printf("Starting service %v...", s.command)

	program := s.command[0]
	argv := s.command[1:]
	cmd := exec.CommandContext(s.ctx, program, argv...)
	s.cmd = cmd

	stdoutRd, err := cmd.StdoutPipe()
	if err != nil {
		eventChan <- err
		return
	}
	go pump(stdoutRd, "Service (out):", os.Stdout, eventChan)

	stderrRd, err := cmd.StderrPipe()
	if err != nil {
		eventChan <- err
		return
	}
	go pump(stderrRd, "Service (err):", os.Stderr, eventChan)

	err = cmd.Start()
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

func (s *Service) Kill() {
	if s.cmd != nil {
		_ = s.cmd.Process.Kill()
	}
}

func (s *Service) Start(eventChan chan interface{}) {
	go s.run(eventChan)
}
