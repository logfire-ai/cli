package alerts_pause

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
	"github.com/spf13/cobra"
)

type PauseAlertOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	AlertPause  bool
	AlertId     []string
}

func NewPauseAlertCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &PauseAlertOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "pause",
		Short: "pause alerts",
		Long:  "pause alerts",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire alerts pause

			# start argument setup
			$ logfire alerts pause --team-id <team-id> --alert-pause <true|false> --alert-id <alert-id> (multiple alerts are allowed)
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PauseAlertRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team id from which alert is to be pause or unpaused.")
	cmd.Flags().BoolVarP(&opts.AlertPause, "alert-pause", "p", false, "Alert pause true or false.")
	cmd.Flags().StringSliceVarP(&opts.AlertId, "alert-id", "a", nil, "alerts to be paused or unpaused. (multiple alerts are allowed)")
	return cmd
}

func PauseAlertRun(opts *PauseAlertOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)
		opts.AlertId, _ = pre_defined_prompters.AskAlertIds(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		if len(opts.AlertId) == 0 {
			fmt.Fprintf(opts.IO.ErrOut, "%s No alerts to pause/unpause\n", cs.FailureIcon())
			os.Exit(0)
		}

		opts.AlertPause, err = opts.Prompter.Confirm("Do you want to pause the alerts? (Yes = Pause, No = Un-pause", false)

	} else {
		if opts.TeamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}

		if opts.AlertId == nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Alerts id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}
	}

	err = APICalls.PauseAlert(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId,
		opts.AlertId, opts.AlertPause)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		if opts.AlertPause == false {
			fmt.Fprintf(opts.IO.Out, "%s Alerts unpaused successfully!\n", cs.SuccessIcon())

		} else {
			fmt.Fprintf(opts.IO.Out, "%s Alerts paused successfully!\n", cs.SuccessIcon())
		}
	}
}
