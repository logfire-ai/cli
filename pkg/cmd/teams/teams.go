package teams

import (
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/teams/member_invite"
	"github.com/logfire-sh/cli/pkg/cmd/teams/member_list"
	"github.com/logfire-sh/cli/pkg/cmd/teams/member_remove"
	"github.com/logfire-sh/cli/pkg/cmd/teams/member_update"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_create"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_delete"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_list"
	"github.com/logfire-sh/cli/pkg/cmd/teams/team_update"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

type PromptTeamsOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	Choice      string
}

var choices = []string{"Create team", "List teams", "Delete team", "Update team", "Invite members", "List members", "Remove member", "Update member"}

func NewCmdTeam(f *cmdutil.Factory) *cobra.Command {
	opts := &PromptTeamsOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:     "teams <command>",
		Short:   "Team details",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PromptIntegrationsRun(opts)

			switch opts.Choice {
			case choices[0]:
				team_create.NewCreateCmd(f).Run(cmd, []string{})
			case choices[1]:
				team_list.NewListCmd(f).Run(cmd, []string{})
			case choices[2]:
				team_delete.NewDeleteCmd(f).Run(cmd, []string{})
			case choices[3]:
				team_update.NewUpdateCmd(f).Run(cmd, []string{})
			case choices[4]:
				member_invite.NewMemberInviteCmd(f).Run(cmd, []string{})
			case choices[5]:
				member_list.NewMemberListCmd(f).Run(cmd, []string{})
			case choices[6]:
				member_remove.NewMemberRemoveCmd(f).Run(cmd, []string{})
			case choices[7]:
				member_update.NewMemberUpdateCmd(f).Run(cmd, []string{})
			}
		},
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

func PromptIntegrationsRun(opts *PromptTeamsOptions) {
	cs := opts.IO.ColorScheme()
	if !opts.Interactive {
		return
	}

	if opts.Interactive {
		err := errors.New("")
		opts.Choice, err = opts.Prompter.Select("What do you want to do?", "", choices)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read choice\n", cs.FailureIcon())
			return
		}
	}
}
