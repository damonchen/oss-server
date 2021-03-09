package main

import (
	"github.com/damonchen/oss-server/internal/config"
	"io"

	"github.com/spf13/cobra"
)

func newRootCmd(cfg *config.Configuration, outer io.Writer, args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "osv",
		Short:        "oss server for aliyun and tencent oss, etc.",
		SilenceUsage: true,
	}

	flags := cmd.PersistentFlags()
	flags.ParseErrorsWhitelist.UnknownFlags = true
	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(
		newWebCmd(cfg, outer, args),
		newVersionCmd(outer),
	)
	return cmd, nil
}
