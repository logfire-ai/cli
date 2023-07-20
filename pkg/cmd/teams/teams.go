package teams

import (
	"github.com/logfire-sh/cli/pkg/cmd/teams/member_invite"
	"github.com/logfire-sh/cli/pkg/cmd/teams/member_list"
	"github.com/logfire-sh/cli/pkg/cmd/teams/member_remove"
	"github.com/logfire-sh/cli/pkg/cmd/teams/member_update"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_create"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_delete"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_list"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_update"
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
	cmd.AddCommand(team_update.NewUpdateCmd(f))
	cmd.AddCommand(team_delete.NewDeleteCmd(f))
	cmd.AddCommand(member_list.NewMemberListCmd(f))
	cmd.AddCommand(member_invite.NewMemberInviteCmd(f))
	cmd.AddCommand(member_remove.NewMemberRemoveCmd(f))
	cmd.AddCommand(member_update.NewMemberUpdateCmd(f))
	return cmd
}
