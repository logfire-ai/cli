package alerts_create

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
		Short: "Create a new alert",
		Long:  "Create a new alert",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire alerts create

			# start argument setup
			$ logfire alerts create --team-id <team-id> --name <name> --view-id <view-id> 
			--number-of-records <0-1000000> --within-seconds <0-10000> --integrations-id <integrations-id> (multiple-integrations-ids supported)
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			CreateAlertRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team id for which alert is to be created.")
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name for the alert.")
	cmd.Flags().StringVarP(&opts.ViewId, "view-id", "v", "", "View id for which alert is to be created.")
	cmd.Flags().Uint32VarP(&opts.NumberOfRecords, "number-of-records", "r", 0, "number of records at when alerts should be triggered.")
	cmd.Flags().Uint32VarP(&opts.WithinSeconds, "within-seconds", "w", 0, "Time range where number of records should occur for alerts to be triggered.")
	cmd.Flags().StringSliceVarP(&opts.IntegrationsId, "integrations-id", "i", nil, "integration to be alerted. (multiple integrations are allowed")
	return cmd
}

func CreateAlertRun(opts *CreateAlertOption) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.TeamId == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
		os.Exit(0)
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

	err = APICalls.CreateAlert(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId,
		opts.Name, opts.ViewId, opts.NumberOfRecords, opts.WithinSeconds, opts.IntegrationsId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Alert created successfully!\n", cs.SuccessIcon())
	}
}
