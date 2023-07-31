package team_create

import (
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type TeamCreateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamName    string
}

func NewCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &TeamCreateOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Args:  cobra.ExactArgs(0),
		Short: "Create teams",
		Long: heredoc.Docf(`
			Create a team.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams create

			# start argument setup
			$ logfire teams create --name <team-name>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			TeamsCreateRun(opts)
		},
	}
	cmd.Flags().StringVar(&opts.TeamName, "name", "", "Name of the team to be created.")
	return cmd
}

func TeamsCreateRun(opts *TeamCreateOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive && opts.TeamName == "" {
		opts.TeamName, err = opts.Prompter.Input("Enter a name for the team:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Name\n", cs.FailureIcon())
			return
		}
	} else {
		if opts.TeamName == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Team name is required.\n", cs.FailureIcon())
		}
	}

	team, err := APICalls.CreateTeam(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamName)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to create team.\n", cs.FailureIcon())
	}

	fmt.Fprintf(opts.IO.Out, "%s Team created successfully.\n%s %s %s %s\n", cs.SuccessIcon(), cs.IntermediateIcon(), team.Name, team.ID, team.Role)
}
