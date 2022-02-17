package service

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func pump(rd io.Reader, prefix string, wr io.Writer, eventChan chan<- interface{}) {

	prefixBytes := []byte(fmt.Sprintf("%s ", prefix))

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

		log.Println(string(prefixBytes), string(line))
		if err != nil {
			eventChan <- err
			return
		}
	}
}
