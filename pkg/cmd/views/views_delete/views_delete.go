package views_delete

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

type ViewsDeleteOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	ViewId      string
}

func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ViewsDeleteOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a view",
		Long:  "Delete a view",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire views delete

			# start argument setup
			$ logfire views delete --team-name <team-name> --view-id <view-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			viewDeleteRun(opts)
		},
	}
	cmd.Flags().StringVar(&opts.TeamId, "team-name", "", "Team name to be deleted.")
	cmd.Flags().StringVar(&opts.ViewId, "view-id", "", "View id to be deleted.")
	return cmd
}

func viewDeleteRun(opts *ViewsDeleteOptions) {
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

	if opts.Interactive {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		opts.ViewId, _ = pre_defined_prompters.AskViewId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
		}

		if opts.ViewId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s View id is required.\n", cs.FailureIcon())
		}
	}

	err = APICalls.DeleteView(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.ViewId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to delete view\n", cs.FailureIcon())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s View deleted successfully\n", cs.SuccessIcon())
	}
}
