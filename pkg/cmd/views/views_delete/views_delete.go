package views_delete

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

type ViewsDeleteOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamID      string
	ViewID      string
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
			$ logfire views delete --team-id <team-id> --view-id <view-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			viewDeleteRun(opts)
		},
	}
	cmd.Flags().StringVar(&opts.TeamID, "team-id", "", "Team id to be deleted.")
	cmd.Flags().StringVar(&opts.ViewID, "view-id", "", "View id to be deleted.")
	return cmd
}

func viewDeleteRun(opts *ViewsDeleteOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.TeamID == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
	}

	if opts.ViewID == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s View id is required.\n", cs.FailureIcon())
	}

	err = APICalls.DeleteView(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamID, opts.ViewID)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to delete view\n", cs.FailureIcon())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s View deleted successfully\n", cs.SuccessIcon())
	}
}
