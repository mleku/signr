package cmd

import (
	"fmt"
	"os"
)

func printErr(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}
