package integrations_create

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

type CreateIntegrationOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive     bool
	TeamId          string
	Name            string
	Type            string
	Description     string
	IntegrationType string
	Id              string
}

func NewCreateIntegrationsCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateIntegrationOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Integration",
		Long:  "Create a new Integration",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire integrations create

			# start argument setup
			$ logfire integrations create --team-id <team-id> --name <name> --description <description>
			  --integration-type <email | webhook | slack> --id <email-id | webhook-id | slack-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			CreateIntegrationRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team id for which Integration is to be created.")
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name for the Integration.")
	cmd.Flags().StringVarP(&opts.Type, "description", "d", "", "Description for the Integration.")
	cmd.Flags().StringVarP(&opts.IntegrationType, "integration-type", "", "", "Type of Integration [email, webhook, slack] (Any one).")
	cmd.Flags().StringVarP(&opts.Id, "id", "i", "", "email-id | webhook-id | slack-id")

	return cmd
}

func CreateIntegrationRun(opts *CreateIntegrationOptions) {
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

	err = APICalls.CreateIntegration(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId,
		opts.Name, opts.Description, opts.Id, opts.IntegrationType)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Integration created successfully!\n", cs.SuccessIcon())
	}
}
