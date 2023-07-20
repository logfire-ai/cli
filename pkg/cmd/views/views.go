package views

import (
	"github.com/logfire-sh/cli/pkg/cmd/views/views_delete"
	"github.com/logfire-sh/cli/pkg/cmd/views/views_list"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdViews(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "views <command>",
		Short:   "Views",
		GroupID: "core",
	}

	cmd.AddCommand(views_delete.NewDeleteCmd(f))
	cmd.AddCommand(views_list.NewViewListCmd(f))
	return cmd
}
