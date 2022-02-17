package service

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func pump(rd io.Reader, prefix string, eventChan chan<- interface{}) {

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

		log.Println(prefixStr, string(line))
		if err != nil {
			eventChan <- err
			return
		}
	}
}
