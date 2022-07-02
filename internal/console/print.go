package console

import (
	"fmt"
	"os"
)

func Infof(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func Info(messages ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, messages...)
}
