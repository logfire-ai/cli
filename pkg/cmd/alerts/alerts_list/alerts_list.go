package alerts_list

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/helpers"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type ListAlertOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
}

func NewListAlertCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListAlertOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all alerts",
		Long:  "list all alerts",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire alerts list

			# start argument setup
			$ logfire alerts list --team-name <team-name>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			ListAlertRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-name", "t", "", "Team name from which alerts are to be listed.")
	return cmd
}

func ListAlertRun(opts *ListAlertOptions) {
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
		if opts.TeamId == "" {
			opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)
		}
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
		}
	}

	data, err := APICalls.ListAlert(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else if len(data) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Alert-Id"})

		for _, i2 := range data {
			table.Append([]string{i2.Name, i2.Id})
		}

		table.Render()
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s No alerts created. Please create an alert\n", cs.FailureIcon())
		os.Exit(0)
	}
}
