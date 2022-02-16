package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceStarts(t *testing.T) {
	ctx := context.Background()
	service := Service{ctx: ctx, command: []string{"bash", "-c", "sleep 1"}}
	eventChan := make(chan interface{})

	go service.Run(eventChan)

	event := <-eventChan
	err, isError := event.(error)
	if isError {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "service started", event)

	event = <-eventChan
	err, isError = event.(error)
	if isError {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "service stopped", event)
}

func TestServiceFailureDetected(t *testing.T) {
	ctx := context.Background()
	service := Service{ctx: ctx, command: []string{"false"}}
	eventChan := make(chan interface{})

	go service.Run(eventChan)

	event := <-eventChan
	err, isError := event.(error)
	if isError {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "service started", event)

	event = <-eventChan
	err, isError = event.(error)
	if !isError {
		t.Fail()
		return
	}

	assert.NotNil(t, err)
}
