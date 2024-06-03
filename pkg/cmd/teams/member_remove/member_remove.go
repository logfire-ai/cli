package member_remove

import (
	"fmt"
	"net/http"
	"os"

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
		Use:   "remove-member",
		Short: "remove member of a team",
		Long: heredoc.Docf(`
			remove member of a team.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams remove-member

			# start argument setup
			$ logfire teams remove-member --team-name <team-name> --memberid <member-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive && opts.TeamId == "" {
				fmt.Fprint(opts.IO.ErrOut, "team-name is required.")
			}

			RemoveMemberRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-name", "", "Team name for which member is to be deleted.")
	cmd.Flags().StringVar(&opts.MemberId, "memberid", "", "Member id of the member to be deleted")
	return cmd
}

func RemoveMemberRun(opts *MemberRemoveOptions) {
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

	if opts.Interactive && opts.TeamId == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		opts.MemberId, _ = pre_defined_prompters.AskMemberId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
		}

		if opts.MemberId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Member id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}
	}

	err = APICalls.RemoveMember(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.MemberId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Member removed successfully!\n", cs.SuccessIcon())
	}
}
