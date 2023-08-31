package integrations_list

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type ListIntegrationsOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
}

func NewListIntegrationsCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListIntegrationsOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all integrations",
		Long:  "List all integrations",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire integrations list

			# start argument setup
			$ logfire integrations list --team-id <team-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			ListIntegrationsRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team id from which integrations are to be listed.")
	return cmd
}

func ListIntegrationsRun(opts *ListIntegrationsOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive {
		if opts.TeamId == "" {
			opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)
		}
	} else {
		if opts.TeamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}

	}

	data, err := APICalls.GetIntegrationsList(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else if len(data) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Integration-Id"})

		for _, i2 := range data {
			table.Append([]string{i2.Name, i2.Id})
		}

		table.Render()
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s No alerts created. Please create an alert\n", cs.FailureIcon())
		os.Exit(0)
	}
}
