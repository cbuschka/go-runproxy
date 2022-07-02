package console

import (
	"fmt"
	"os"
)

var verbose = false

func EnableVerboseOutput() {
	verbose = true
}

func Infof(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	Info(s)
}

func Info(messages ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, messages...)
}

func Debugf(format string, args ...interface{}) {
	if !verbose {
		return
	}

	s := fmt.Sprintf(format, args...)
	Info(s)
}

func Debug(messages ...interface{}) {
	if !verbose {
		return
	}

	Info(messages...)
}
