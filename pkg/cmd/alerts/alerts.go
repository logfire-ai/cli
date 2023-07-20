package alerts

import (
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_create"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_delete"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_list"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_pause"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_update"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAlerts(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alerts <command>",
		Short:   "alerts",
		GroupID: "core",
	}

	cmd.AddCommand(alerts_create.NewCreateAlertCmd(f))
	cmd.AddCommand(alerts_list.NewListAlertCmd(f))
	cmd.AddCommand(alerts_delete.NewDeleteAlertCmd(f))
	cmd.AddCommand(alerts_pause.NewPauseAlertCmd(f))
	cmd.AddCommand(alerts_update.NewAlertUpdateCmd(f))
	return cmd
}
