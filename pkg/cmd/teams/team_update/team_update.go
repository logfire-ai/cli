package team_update

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

type TeamUpdateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	TeamName    string
}

func NewUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &TeamUpdateOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update a team",
		Long:  "update a team",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams update

			# start argument setup
			$ logfire teams update --teamid <team-id> --name <new-name>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive && opts.TeamName == "" {
				fmt.Fprint(opts.IO.ErrOut, "new name is required.\n")
			}

			if !opts.Interactive && opts.TeamId == "" {
				fmt.Fprint(opts.IO.ErrOut, "team id is required.\n")
			}

			teamUpdateRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamName, "name", "", "new team name to be updated.")
	cmd.Flags().StringVar(&opts.TeamId, "teamid", "", "Team id to be updated.")
	return cmd
}

func teamUpdateRun(opts *TeamUpdateOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive && opts.TeamId == "" && opts.TeamName == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		opts.TeamName, err = opts.Prompter.Input("Enter a new name for the team:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Name\n", cs.FailureIcon())
			return
		}
	} else {
		if opts.TeamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
		}

		if opts.TeamName == "" {
			fmt.Fprint(opts.IO.ErrOut, "new name is required.\n", cs.FailureIcon())
		}
	}

	resp, err := APICalls.UpdateTeam(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.TeamName)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to update team\n", cs.FailureIcon())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Team updated successfully, Name: %s, ID: %s\n", cs.SuccessIcon(), resp.Name, resp.ID)
	}
}
