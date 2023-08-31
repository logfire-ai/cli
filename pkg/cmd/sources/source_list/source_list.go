package source_list

import (
	"fmt"
	"net/http"
	"os"

	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/olekukonko/tablewriter"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type SourceListOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	TeamId string
}

func NewSourceListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SourceListOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Args:  cobra.ExactArgs(0),
		Short: "Get sources",
		Long: heredoc.Docf(`
			Get sources for a particular team.

			The user is prompted with the teams. User can select a team to show the sources.
		`),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sources list

			# start argument setup
			$ logfire sources list --team-id <team-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if opts.TeamId == "" && !opts.Interactive {
				fmt.Fprint(opts.IO.ErrOut, "team-id is required.\n")
				return
			}

			sourceListRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-id", "", "Team ID for which the sources will be fetched.")
	return cmd
}

func sourceListRun(opts *SourceListOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Interactive && opts.TeamId == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)
	} else {
		if opts.TeamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s team-id is required.\n", cs.FailureIcon())
			return
		}
	}

	sources, err := APICalls.GetAllSources(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		return
	} else if len(sources) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Source-Id", "Platform", "Token"})

		for _, i2 := range sources {
			table.Append([]string{i2.Name, i2.ID, i2.Platform, i2.SourceToken})
		}

		table.Render()
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s No sources created. Please create a source\n", cs.FailureIcon())
		os.Exit(0)
	}
}
