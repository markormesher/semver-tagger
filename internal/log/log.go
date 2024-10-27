package log

import (
	"fmt"
	"os"
	"strings"
)

var Verbose = false

func Debug(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}

	if Verbose {
		fmt.Printf("DEBUG: "+format, args...)
	}
}

func Info(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}

	fmt.Printf(format, args...)
}

func Warn(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}

	fmt.Fprintf(os.Stderr, "WARN: "+format, args...)
}

func Error(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}

	fmt.Fprintf(os.Stderr, "ERROR: "+format, args...)
}
