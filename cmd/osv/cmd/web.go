package cmd

import (
	"io"

	"github.com/damonchen/oss-server/cmd/osv/cmd/require"
	"github.com/damonchen/oss-server/internal/action"
	"github.com/damonchen/oss-server/internal/config"
	"github.com/spf13/cobra"
)

func newWebCmd(cfg *config.Configuration, outer io.Writer, args []string) *cobra.Command {
	server := action.Web{}
	cmd := &cobra.Command{
		Use:   "web",
		Short: "web server",
		Long:  "",
		Args:  require.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := server.Run()
			return err
		},
	}

	f := cmd.Flags()
	f.StringVar(&server.CfgFile, "cfg", "", "config filename path")

	return cmd
}
