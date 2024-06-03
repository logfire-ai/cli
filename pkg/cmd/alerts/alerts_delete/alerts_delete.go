package alerts_delete

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
	"github.com/spf13/cobra"
)

type DeleteAlertOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	AlertId     []string
}

func NewDeleteAlertCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteAlertOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete alerts",
		Long:  "delete alerts",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire alerts delete

			# start argument setup
			$ logfire alerts delete --team-name <team-name> --alert-id <alert-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			DeleteAlertRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-name", "t", "", "Team name from which alert is to be deleted.")
	cmd.Flags().StringSliceVarP(&opts.AlertId, "alert-id", "a", nil, "alerts to be deleted. (multiple alerts are allowed)")
	return cmd
}

func DeleteAlertRun(opts *DeleteAlertOptions) {
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

		opts.AlertId, _ = pre_defined_prompters.AskAlertIds(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
		}

		if opts.AlertId == nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Alerts id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}

	}

	err = APICalls.DeleteAlert(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId,
		opts.AlertId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Alerts deleted successfully!\n", cs.SuccessIcon())
	}
}
