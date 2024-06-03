package integrations_create

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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

type CreateIntegrationOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive     bool
	TeamId          string
	Name            string
	Description     string
	IntegrationType string
	Id              string
}

var IntegrationOptions = []string{
	"Email", "Slack", "Webhook", "Member",
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
			$ logfire integrations create --team-name <team-name> --name <name> --description <description>
			  --integration-type <email | webhook | slack> --id <email-id | webhook-id | slack-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			CreateIntegrationRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-name", "t", "", "Team name for which Integration is to be created.")
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name for the Integration.")
	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Description for the Integration.")
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

	client := http.Client{}

	if opts.TeamId != "" {
		teamId := helpers.TeamNameToTeamId(&client, cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		if teamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s no team with name: %s found.\n", cs.FailureIcon(), opts.TeamId)
			return
		}

		opts.TeamId = teamId
	}

	if opts.Interactive && opts.TeamId == "" && opts.Name == "" && opts.Description == "" && opts.IntegrationType == "" && opts.Id == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		opts.Name, err = opts.Prompter.Input("Enter a name for the alert:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Name\n", cs.FailureIcon())
			return
		}

		opts.Description, err = opts.Prompter.Input("Enter a Description for the alert:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Name\n", cs.FailureIcon())
			return
		}

		opts.IntegrationType, err = opts.Prompter.Select("Select a type of the Integration:", "", IntegrationOptions)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Integration type\n", cs.FailureIcon())
			return
		}

		opts.Id, err = opts.Prompter.Input("Enter ID of the Integration", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Integration id\n", cs.FailureIcon())
			return
		}
	} else {
		if opts.TeamId == "" {
			opts.TeamId = cfg.Get().TeamId
		}

		if opts.Name == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Name is required.\n", cs.FailureIcon())
			os.Exit(0)
		}

		if opts.IntegrationType == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Integration type is required.\n", cs.FailureIcon())
			os.Exit(0)
		}

		if opts.Id == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Integration id is required.\n", cs.FailureIcon())
			os.Exit(0)
		}
	}

	err = APICalls.CreateIntegration(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId,
		opts.Name, opts.Description, opts.Id, strings.ToLower(opts.IntegrationType))
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Integration created successfully!\n", cs.SuccessIcon())
	}
}
