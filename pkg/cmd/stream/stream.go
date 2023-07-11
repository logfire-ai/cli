package stream

import (
	"github.com/logfire-sh/cli/pkg/cmd/stream/livetail"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdStream(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stream",
		Short:   "Logs streaming",
		GroupID: "core",
	}

	cmd.AddCommand(livetail.NewLivetailCmd(f))

	return cmd
}
