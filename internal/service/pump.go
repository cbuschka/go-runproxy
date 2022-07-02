package service

import (
	"bufio"
	"github.com/cbuschka/go-runproxy/internal/console"
	"io"
	"regexp"
)

func pump(rd io.Reader, startupMessageMatchPattern *regexp.Regexp, eventChan chan<- interface{}) {

	bufferedRd := bufio.NewReader(rd)

	startupMessageSeen := false
	for {
		line, _, err := bufferedRd.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			eventChan <- err
			return
		}

		if len(line) == 0 {
			continue
		}

		lineStr := string(line)
		if !startupMessageSeen && startupMessageMatchPattern != nil && startupMessageMatchPattern.MatchString(lineStr) {
			eventChan <- "startup message seen"
			startupMessageSeen = true
		}

		console.Info(lineStr)
		if err != nil {
			eventChan <- err
			return
		}
	}
}
