package integrations

import (
	"github.com/logfire-sh/cli/pkg/cmd/integrations/integrations_create"
	"github.com/logfire-sh/cli/pkg/cmd/integrations/integrations_delete"
	"github.com/logfire-sh/cli/pkg/cmd/integrations/integrations_list"
	"github.com/logfire-sh/cli/pkg/cmd/integrations/integrations_update"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdIntegrations(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "integrations <command>",
		Short:   "integrations",
		GroupID: "core",
	}

	cmd.AddCommand(integrations_create.NewCreateIntegrationsCmd(f))
	cmd.AddCommand(integrations_list.NewListIntegrationsCmd(f))
	cmd.AddCommand(integrations_delete.NewDeleteIntegrationCmd(f))
	cmd.AddCommand(integrations_update.NewUpdateIntegrationsCmd(f))
	return cmd
}
