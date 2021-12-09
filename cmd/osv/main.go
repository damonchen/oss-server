package main

import (
	"fmt"
	"os"

	"github.com/damonchen/oss-server/cmd/osv/cmd"
	"github.com/damonchen/oss-server/internal/config"
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
	cmd, err := cmd.NewRootCmd(cfg, os.Stdout, os.Args[1:])
	if err != nil {
		warning("%+v", err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		debug("%+v", err)
		os.Exit(1)
	}
}
