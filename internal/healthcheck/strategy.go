package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

type HealthcheckStrategy interface {
	IsAvailable(ctx context.Context) (bool, error)
	String() string
}

type HttpHealthcheckStrategy struct {
	endpointAddress string
}

type CommandHealthcheckStrategy struct {
	command []string
}

type DummyHealthcheckStrategy struct {
}

func (p *DummyHealthcheckStrategy) IsAvailable(ctx context.Context) (bool, error) {
	return true, nil
}

func (p *HttpHealthcheckStrategy) IsAvailable(ctx context.Context) (bool, error) {
	endpointUrl, err := url.Parse(p.endpointAddress)
	if err != nil {
		return false, err
	}

	client := &http.Client{}

	req := http.Request{
		URL:    endpointUrl,
		Method: "GET",
	}

	resp, err := client.Do(&req)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}

func (p *CommandHealthcheckStrategy) IsAvailable(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, p.command[0], p.command[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return false, err
	}

	err = cmd.Wait()
	if err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		if exitCode != 0 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (p *CommandHealthcheckStrategy) String() string {
	return fmt.Sprintf("CommandHealthcheckStrategy{command=%v}", p.command)
}

func (p *HttpHealthcheckStrategy) String() string {
	return fmt.Sprintf("HttpHealthcheckStrategy{endpointAddress=%s}", p.endpointAddress)
}

func (p *DummyHealthcheckStrategy) String() string {
	return "dummy"
}
