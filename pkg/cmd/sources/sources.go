package sources

import (
	"github.com/logfire-sh/cli/pkg/cmd/sources/source_create"
	"github.com/logfire-sh/cli/pkg/cmd/sources/source_delete"
	"github.com/logfire-sh/cli/pkg/cmd/sources/source_list"
	"github.com/logfire-sh/cli/pkg/cmd/sources/source_update"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSource(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sources <command>",
		Short:   "Get source for a team",
		GroupID: "core",
	}

	cmd.AddCommand(source_list.NewSourceListCmd(f))
	cmd.AddCommand(source_create.NewSourceCreateCmd(f))
	cmd.AddCommand(source_update.NewSourceUpdateCmd(f))
	cmd.AddCommand(source_delete.NewSourceDeleteCmd(f))
	return cmd
}
