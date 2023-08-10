package member_list

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type MemberListOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
}

func NewMemberListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &MemberListOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "list-members",
		Args:  cobra.ExactArgs(0),
		Short: "List team members",
		Long: heredoc.Docf(`
			List all the members of a team.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams list-members

			# start argument setup
			$ logfire teams invite-members --teamid <team-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive && opts.TeamId == "" {
				fmt.Fprint(opts.IO.ErrOut, "team-id is required.")
			}

			listMembersRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-id", "", "Team id for which members are to be fetched.")
	return cmd
}

func listMembersRun(opts *MemberListOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive && opts.TeamId == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)
	} else {
		if opts.TeamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}
	}

	members, err := APICalls.MembersList(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"First name", "Last name", "Profile-Id", "Role"})

		for _, i2 := range members {
			table.Append([]string{i2.FirstName, i2.LastName, i2.ProfileId, i2.Role.String()})
		}

		table.Render()
	}
}
