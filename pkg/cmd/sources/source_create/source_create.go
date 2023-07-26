package source_create

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

var platformOptions = make([]string, 0, len(models.PlatformMap))

type SourceCreateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	TeamId     string
	SourceName string
	Platform   string
}

func NewSourceCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SourceCreateOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Args:  cobra.ExactArgs(0),
		Short: "Create source",
		Long: heredoc.Docf(`
			Create a source for the particular team.
		`),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sources create

			# start argument setup
			$ logfire sources create --teamid <team-id> --name <source-name> --platform <platform>
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

				if opts.Platform == "" {
					fmt.Fprint(opts.IO.ErrOut, "platform is required.\n")
					return
				}
			}

			sourceCreateRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-id", "", "Team ID for which the source will be created.")
	cmd.Flags().StringVar(&opts.SourceName, "name", "", "Name of the source to be created.")
	cmd.Flags().StringVar(&opts.Platform, "platform", "", "Platform name for which you want to create source.")
	return cmd
}

func sourceCreateRun(opts *SourceCreateOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Interactive {
		opts.TeamId, err = opts.Prompter.Input("Enter TeamId:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read TeamId\n", cs.FailureIcon())
			return
		}

		opts.SourceName, err = opts.Prompter.Input("Enter Source name:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Source name\n", cs.FailureIcon())
			return
		}

		opts.Platform, err = opts.Prompter.Select("Enter Platform name:", "", platformOptions)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Platform name\n", cs.FailureIcon())
			return
		}
	}

	if opts.TeamId == "" || opts.SourceName == "" || opts.Platform == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s team-id, name and plaform are required.\n", cs.FailureIcon())
		return
	}
	source, err := APICalls.CreateSource(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.SourceName, opts.Platform)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		return
	}

	fmt.Fprintf(opts.IO.Out, "%s Successfully created source for team-id %s\n", cs.SuccessIcon(), opts.TeamId)
	fmt.Fprintf(opts.IO.Out, "%s %s %s %s %s\n", cs.IntermediateIcon(), source.Name, source.ID, source.SourceToken, source.Platform)
}
