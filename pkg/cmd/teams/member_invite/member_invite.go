package member_invite

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/helpers"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type MemberInviteOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	Email       []string
}

func NewMemberInviteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &MemberInviteOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "invite-members",
		Short: "invite-members to a team",
		Long: heredoc.Docf(`
			invite-members to a team.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams invite-members

			# start argument setup
			$ logfire teams invite-members --team-name <team-name> --email <email> --email <email> (multiple values supported)
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			InviteMembersRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-name", "", "Team name for which members are to be invited.")
	cmd.Flags().StringSliceVarP(&opts.Email, "email", "e", nil, "Email addresses (multiple values supported).")
	return cmd
}

func InviteMembersRun(opts *MemberInviteOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	client := http.Client{}

	if opts.TeamId != "" {
		teamId := helpers.TeamNameToTeamId(&client, cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		if teamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s no team with name: %s found.\n", cs.FailureIcon(), opts.TeamId)
			return
		}

		opts.TeamId = teamId
	}

	if opts.Interactive && opts.TeamId == "" && opts.Email == nil {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		emails, err := opts.Prompter.Input("Enter email addresses to invite (Multiple email can be entered separated by a comma).", "")
		if err != nil {
			return
		}

		opts.Email = strings.Split(emails, ",")
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
		}

		if opts.Email == nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Atleast one email address is required.\n", cs.FailureIcon())
			os.Exit(0)
		}
	}

	err = APICalls.InviteMembers(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.Email)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Members invited successfully!\n", cs.SuccessIcon())
	}
}
