package service

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
)

func pump(rd io.Reader, prefix string, startupMessageMatchPattern *regexp.Regexp, eventChan chan<- interface{}) {

	prefixStr := fmt.Sprintf("%s ", prefix)

	bufferedRd := bufio.NewReader(rd)

	for {
		line, _, err := bufferedRd.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			eventChan <- err
			return
		}

		lineStr := string(line)
		if startupMessageMatchPattern != nil && startupMessageMatchPattern.MatchString(lineStr) {
			eventChan <- "startup message seen"
		}

		log.Println(prefixStr, lineStr)
		if err != nil {
			eventChan <- err
			return
		}
	}
}
