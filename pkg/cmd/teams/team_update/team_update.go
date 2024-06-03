package team_update

import (
	"fmt"
	"net/http"

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
			$ logfire teams update --team-name <team-name> --name <new-name>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := opts.Config()
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "Failed to read config\n")
			}

			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive && opts.TeamName == "" {
				fmt.Fprint(opts.IO.ErrOut, "new name is required.\n")
			}

			if !opts.Interactive && opts.TeamId == "" {
				opts.TeamId = cfg.Get().TeamId
			}

			teamUpdateRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamName, "name", "", "new team name to be updated.")
	cmd.Flags().StringVar(&opts.TeamId, "team-name", "", "Team name to be updated.")
	return cmd
}

func teamUpdateRun(opts *TeamUpdateOptions) {
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

	if opts.Interactive && opts.TeamId == "" && opts.TeamName == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		opts.TeamName, err = opts.Prompter.Input("Enter a new name for the team:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Name\n", cs.FailureIcon())
			return
		}
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
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
