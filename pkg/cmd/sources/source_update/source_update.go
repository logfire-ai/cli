package source_update

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

type SourceUpdateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	TeamId     string
	SourceId   string
	SourceName string
}

func NewSourceUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SourceUpdateOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "update",
		Args:  cobra.ExactArgs(0),
		Short: "update source",
		Long: heredoc.Docf(`
			update a source for the particular team.
		`),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sources update

			# start argument setup
			$ logfire sources update --teamid <team-id> --sourceid <source-id> --name <new-name>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive {
				if opts.TeamId == "" {
					fmt.Fprint(opts.IO.ErrOut, "team-id is required.\n")
					return
				}

				if opts.SourceName == "" {
					fmt.Fprint(opts.IO.ErrOut, "name is required.\n")
					return
				}

				if opts.SourceId == "" {
					fmt.Fprint(opts.IO.ErrOut, "source-id is required.\n")
					return
				}
			}

			sourceUpdateRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "teamid", "", "Team ID for which the source will be updated.")
	cmd.Flags().StringVar(&opts.SourceId, "sourceid", "", "Source ID for which the source will be updated.")
	cmd.Flags().StringVar(&opts.SourceName, "name", "", "New name for the source to be updated.")
	return cmd
}

func sourceUpdateRun(opts *SourceUpdateOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Interactive {
	}

	if opts.TeamId == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
		return
	}

	if opts.SourceName == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s New name is required.\n", cs.FailureIcon())
		return
	}

	if opts.SourceId == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Source id is required.\n", cs.FailureIcon())
		return
	}

	source, err := APICalls.UpdateSource(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.SourceId, opts.SourceName)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		return
	}

	fmt.Fprintf(opts.IO.Out, "%s Successfully update source for source-id %s\n", cs.SuccessIcon(), opts.SourceId)
	fmt.Fprintf(opts.IO.Out, "%s %s %s %s %s\n", cs.IntermediateIcon(), source.Name, source.ID, source.SourceToken, source.Platform)
}
