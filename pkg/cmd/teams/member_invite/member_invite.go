package member_invite

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
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

func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &MemberInviteOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "invite members",
		Short: "invite members to a team",
		Long: heredoc.Docf(`
			invite members to a team.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams invite members

			# start argument setup
			$ logfire teams invite members --teamid <team-id> --email <email> --email <email> (multiple values supported)
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive && opts.TeamId == "" {
				fmt.Fprint(opts.IO.ErrOut, "team-id is required.")
			}

			InviteMembersRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "teamid", "", "Team id for which members are to be fetched.")
	cmd.Flags().StringSliceVarP(&opts.Email, "email", "e", nil, "Email addresses (multiple values supported).")
	return cmd
}

func InviteMembersRun(opts *MemberInviteOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	err = APICalls.InviteMembers(opts.HttpClient(), cfg.Get().Token, opts.TeamId, opts.Email)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Members invited successfully!\n", cs.SuccessIcon())
	}
}
