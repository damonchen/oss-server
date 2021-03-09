package main

import (
	"fmt"
	"github.com/damonchen/oss-server/cmd/osv/require"
	"github.com/damonchen/oss-server/internal/version"
	"github.com/spf13/cobra"
	"io"
)

type versionOptions struct {
	short bool
}

func newVersionCmd(out io.Writer) *cobra.Command {
	o := versionOptions{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the version",
		Args:  require.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&o.short, "short", false, "print the version number")

	return cmd
}

func (o *versionOptions) run(out io.Writer) error {
	_, _ = fmt.Fprintln(out, formatVersion(o.short))
	return nil
}

func formatVersion(short bool) string {
	v := version.Get()
	if short {
		if len(v.GitCommit) >= 7 {
			return fmt.Sprintf("%s+g%s", v.Version, v.GitCommit[:7])
		}
		return version.GetVersion()
	}
	return fmt.Sprintf("%#v", v)
}
