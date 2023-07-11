package teams

import (
	"github.com/logfire-sh/cli/pkg/cmd/teams/member_list"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_create"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_list"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdTeam(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "teams <command>",
		Short:   "Team details",
		GroupID: "core",
	}

	cmd.AddCommand(team_list.NewListCmd(f))
	cmd.AddCommand(team_create.NewCreateCmd(f))
	cmd.AddCommand(member_list.NewListCmd(f))
	return cmd
}
