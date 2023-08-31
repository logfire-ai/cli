package team_list

import (
	"fmt"
	"net/http"
	"os"

	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/olekukonko/tablewriter"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type TeamOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
}

func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &TeamOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Args:  cobra.ExactArgs(0),
		Short: "List all the teams",
		Long: heredoc.Docf(`
			List teams.
		`, "`"),
		Example: heredoc.Doc(`
			# List all the teams
			$ logfire teams list
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}
			listRun(opts)
		},
	}

	return cmd
}

func listRun(opts *TeamOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	teams, err := APICalls.ListTeams(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else if len(teams) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Team-Id"})

		for _, i2 := range teams {
			table.Append([]string{i2.Name, i2.ID})
		}

		table.Render()
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s No teams created. Please create a team\n", cs.FailureIcon())
		os.Exit(0)
	}
}
