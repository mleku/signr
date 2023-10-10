package cmd

import (
	"fmt"
	"os"
)

func PrintErr(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}
