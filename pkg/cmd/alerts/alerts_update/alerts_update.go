package alerts_update

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

type AlertUpdateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive     bool
	TeamId          string
	AlertId         string
	Name            string
	ViewId          string
	NumberOfRecords uint32
	WithinSeconds   uint32
	IntegrationsId  []string
}

func NewAlertUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &AlertUpdateOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update an alert",
		Long:  "update an alert",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire alerts update

			# start argument setup
			$ logfire alerts update --team-id <team-id> --name <name> --view-id <view-id> 
			--number-of-records <0-1000000> --within-seconds <0-10000> --integrations-id <integrations-id> (multiple-integrations-ids supported)
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			UpdateMemberRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team id for which alert is to be created.")
	cmd.Flags().StringVarP(&opts.AlertId, "alert-id", "a", "", "Team id for which alert is to be created.")
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name for the alert.")
	cmd.Flags().StringVarP(&opts.ViewId, "view-id", "v", "", "View id for which alert is to be created.")
	cmd.Flags().Uint32VarP(&opts.NumberOfRecords, "number-of-records", "r", 0, "number of records at when alerts should be triggered.")
	cmd.Flags().Uint32VarP(&opts.WithinSeconds, "within-seconds", "w", 0, "Time range where number of records should occur for alerts to be triggered.")
	cmd.Flags().StringSliceVarP(&opts.IntegrationsId, "integrations-id", "i", nil, "integration to be alerted. (multiple integrations are allowed")
	return cmd
}

func UpdateMemberRun(opts *AlertUpdateOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		opts.AlertId, _ = pre_defined_prompters.AskAlertId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		updateName, _ := opts.Prompter.Confirm(fmt.Sprintf("Do you want to update the alert name?"), false)
		if updateName {
			opts.Name, _ = opts.Prompter.Input("Enter a new name for the alert:", "")
		}

		updateView, _ := opts.Prompter.Confirm(fmt.Sprintf("Do you want to update the view?"), false)
		if updateView {
			opts.ViewId, _ = pre_defined_prompters.AskViewId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)
		}

		updateNOR, _ := opts.Prompter.Confirm(fmt.Sprintf("Do you want to update the number of records?"), false)
		if updateNOR {
			nor, _ := opts.Prompter.InputInt("number of records at when alerts should be triggered.", 0)
			opts.NumberOfRecords = uint32(nor)
		}

		updateWS, _ := opts.Prompter.Confirm(fmt.Sprintf("Do you want to update the Within-Seconds?"), false)
		if updateWS {
			ws, _ := opts.Prompter.InputInt("Time range where number of records should occur for alerts to be triggered.", 0)
			opts.WithinSeconds = uint32(ws)
		}

		updateIntegrations, _ := opts.Prompter.Confirm(fmt.Sprintf("Do you want to update the integrations?"), false)
		if updateIntegrations {
			opts.IntegrationsId, _ = pre_defined_prompters.AskAlertIntegrationIds(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)
		}
	} else {
		if opts.TeamId == "" {
			fmt.Fprint(opts.IO.ErrOut, "team-id is required.")
			os.Exit(0)
		}

		if opts.AlertId == "" {
			fmt.Fprint(opts.IO.ErrOut, "Alert-id is required.")
			os.Exit(0)
		}
	}

	err = APICalls.UpdateAlert(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId,
		opts.Name, opts.ViewId, opts.NumberOfRecords, opts.WithinSeconds, opts.IntegrationsId, opts.AlertId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Alert updated successfully!\n", cs.SuccessIcon())
	}
}
