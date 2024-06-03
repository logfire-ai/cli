package alerts_create

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

type CreateAlertOption struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive     bool
	TeamId          string
	Name            string
	ViewId          string
	NumberOfRecords uint32
	WithinSeconds   uint32
	IntegrationsId  []string
}

func NewCreateAlertCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateAlertOption{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a new alert",
		Long:  "create a new alert",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire alerts create

			# start argument setup
			$ logfire alerts create --team-name <team-name> --name <name> --view-id <view-id> 
			--number-of-records <0-1000000> --within-seconds <0-10000> --integrations-id <integrations-id> (multiple-integrations-ids supported)
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			CreateAlertRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-name", "t", "", "Team name for which alert is to be created.")
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name for the alert.")
	cmd.Flags().StringVarP(&opts.ViewId, "view-id", "v", "", "View id for which alert is to be created.")
	cmd.Flags().Uint32VarP(&opts.NumberOfRecords, "number-of-records", "r", 0, "number of records at when alerts should be triggered.")
	cmd.Flags().Uint32VarP(&opts.WithinSeconds, "within-seconds", "w", 0, "Time range where number of records should occur for alerts to be triggered.")
	cmd.Flags().StringSliceVarP(&opts.IntegrationsId, "integrations-id", "i", nil, "integration to be alerted. (multiple integrations are allowed)")
	return cmd
}

func CreateAlertRun(opts *CreateAlertOption) {
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
		if opts.TeamId == "" && opts.Name == "" && opts.ViewId == "" && opts.NumberOfRecords == 0 && opts.WithinSeconds == 0 && opts.IntegrationsId == nil {
			opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

			opts.Name, err = opts.Prompter.Input("Enter a name for the alert:", "")
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Name\n", cs.FailureIcon())
				return
			}

			opts.ViewId, _ = pre_defined_prompters.AskViewId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

			nor, err := opts.Prompter.InputInt("Enter a Number at when alerts should be triggered.", 0)
			opts.NumberOfRecords = uint32(nor)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Number Of Records\n", cs.FailureIcon())
				return
			}

			ws, err := opts.Prompter.InputInt("Time range where number of records should occur for alerts to be triggered.", 0)
			opts.WithinSeconds = uint32(ws)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Within Seconds\n", cs.FailureIcon())
				return
			}

			opts.IntegrationsId, _ = pre_defined_prompters.AskAlertIntegrationIds(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		}
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
		}

		if opts.Name == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Name is required.\n", cs.FailureIcon())
			os.Exit(0)
		}

		if opts.ViewId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s View id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}

		if opts.NumberOfRecords == 0 {
			fmt.Fprintf(opts.IO.ErrOut, "%s NumberOfRecords is required.\n", cs.FailureIcon())
			os.Exit(0)
		}

		if opts.WithinSeconds == 0 {
			fmt.Fprintf(opts.IO.ErrOut, "%s WithinSeconds is required.\n", cs.FailureIcon())
			os.Exit(0)
		}

		if opts.IntegrationsId == nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Integrations id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}
	}

	err = APICalls.CreateAlert(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId,
		opts.Name, opts.ViewId, opts.NumberOfRecords, opts.WithinSeconds, opts.IntegrationsId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Alert created successfully!\n", cs.SuccessIcon())
	}
}
