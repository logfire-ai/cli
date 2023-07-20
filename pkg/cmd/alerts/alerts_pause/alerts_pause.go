package alerts_pause

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
	"os"
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
			$ logfire alerts pause --team-id <team-id> --alert-pause <true|false> --alert-id <alert-id>
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

	if opts.TeamId == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
		os.Exit(0)
	}

	if opts.AlertId == nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Alerts id is required.\n", cs.FailureIcon())
		os.Exit(0)
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
