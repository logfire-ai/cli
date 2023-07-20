package member_remove

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

type MemberRemoveOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	MemberId    string
}

func NewMemberRemoveCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &MemberRemoveOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "remove member",
		Short: "remove member of a team",
		Long: heredoc.Docf(`
			remove member of a team.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams remove member

			# start argument setup
			$ logfire teams remove member --teamid <team-id> --memberid <member-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive && opts.TeamId == "" {
				fmt.Fprint(opts.IO.ErrOut, "team-id is required.")
			}

			RemoveMemberRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "teamid", "", "Team id for which member is to be deleted.")
	cmd.Flags().StringVar(&opts.MemberId, "memberid", "", "Member id of the member to be deleted")
	return cmd
}

func RemoveMemberRun(opts *MemberRemoveOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	err = APICalls.RemoveMember(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.MemberId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Member removed successfully!\n", cs.SuccessIcon())
	}
}
