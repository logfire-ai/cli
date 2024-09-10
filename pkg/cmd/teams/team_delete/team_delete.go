package team_delete

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

type TeamDeleteOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
}

func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &TeamDeleteOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a team",
		Long:  "Delete a team",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams delete

			# start argument setup
			$ logfire teams delete --name <team-name>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := opts.Config()
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "Failed to read config\n")
			}

			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive && opts.TeamId == "" {
				opts.TeamId = cfg.Get().TeamId
			}

			teamDeleteRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "name", "n", "", "Team name to be deleted.")
	return cmd
}

func teamDeleteRun(opts *TeamDeleteOptions) {
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
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
		}
	}

	err = APICalls.DeleteTeam(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to delete team\n", cs.FailureIcon())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Team deleted successfully\n", cs.SuccessIcon())
	}
}
