package main

import (
	"fmt"
	"github.com/damonchen/oss-server/internal/config"
	"os"
)

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[DEBUG]: %s\n", format)
	_, _ = fmt.Fprintf(os.Stderr, format, v...)
}

func warning(format string, v ...interface{}) {
	format = fmt.Sprintf("[WARNING]: %s\n", format)
	_, _ = fmt.Fprintf(os.Stderr, format, v...)
}

func main() {
	cfg := new(config.Configuration)
	cmd, err := newRootCmd(cfg, os.Stdout, os.Args[1:])
	if err != nil {
		warning("%+v", err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		debug("%+v", err)
		os.Exit(1)
	}
}
