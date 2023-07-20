package team_delete

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

type TeamDeleteOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamID      string
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
			$ logfire teams delete --teamid <team-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive && opts.TeamID == "" {
				fmt.Fprint(opts.IO.ErrOut, "team id is required.\n")
			}

			teamDeleteRun(opts)
		},
	}
	cmd.Flags().StringVar(&opts.TeamID, "teamid", "", "Team id to be deleted.")
	return cmd
}

func teamDeleteRun(opts *TeamDeleteOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.TeamID == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
	}

	err = APICalls.DeleteTeam(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamID)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to delete team\n", cs.FailureIcon())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Team deleted successfully\n", cs.SuccessIcon())
	}
}
