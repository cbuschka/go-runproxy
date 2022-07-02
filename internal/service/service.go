package service

import (
	"context"
	"fmt"
	"github.com/cbuschka/go-runproxy/internal/console"
	"os/exec"
	"regexp"
)

type Service struct {
	ctx                        context.Context
	command                    []string
	startupMessageMatchPattern *regexp.Regexp

	cmd *exec.Cmd
}

func NewService(ctx context.Context, command []string, startupMessageMatch string) (*Service, error) {
	var startupMessageMatchPattern *regexp.Regexp = nil
	if startupMessageMatch != "" {
		var err error
		startupMessageMatchPattern, err = regexp.Compile(startupMessageMatch)
		if err != nil {
			return nil, err
		}
	}

	return &Service{ctx: ctx, command: command, startupMessageMatchPattern: startupMessageMatchPattern, cmd: nil}, nil
}

func (s *Service) run(eventChan chan<- interface{}) {
	console.Infof("Starting service %v...", s.command)

	program := s.command[0]
	argv := s.command[1:]
	cmd := exec.CommandContext(s.ctx, program, argv...)
	s.cmd = cmd

	stdoutRd, err := cmd.StdoutPipe()
	if err != nil {
		eventChan <- err
		return
	}
	go pump(stdoutRd, s.startupMessageMatchPattern, eventChan)

	stderrRd, err := cmd.StderrPipe()
	if err != nil {
		eventChan <- err
		return
	}
	go pump(stderrRd, s.startupMessageMatchPattern, eventChan)

	err = cmd.Start()
	if err != nil {
		eventChan <- err
		return
	}
	console.Infof("Service %v started.", s.command)

	eventChan <- "service started"

	err = cmd.Wait()
	if err != nil {
		console.Infof("Waiting for service %v failed: %v", s.command, err)
		eventChan <- err
		return
	}

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		console.Infof("Service %v stopped with %d.", s.command, exitCode)
		eventChan <- fmt.Errorf("process exited with %d", exitCode)
		return
	}

	console.Infof("Service %v exited normally.", s.command)

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
