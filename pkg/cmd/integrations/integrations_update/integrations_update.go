package integrations_update

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

type UpdateIntegrationOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive   bool
	TeamId        string
	IntegrationId string
	Name          string
	Description   string
	Id            string
}

func NewUpdateIntegrationsCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UpdateIntegrationOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update an Integration",
		Long:  "update an Integration",
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire integrations update

			# start argument setup
			$ logfire integrations update --team-id <team-id> --integration-id <integration-id> 
				--name <name> --description <description> --id <email-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			UpdateIntegrationRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team id for which Integration is to be updated.")
	cmd.Flags().StringVarP(&opts.IntegrationId, "integration-id", "", "", "Integration id for which settings are to be updated.")
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name for the Integration.")
	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Description for the Integration.")
	cmd.Flags().StringVarP(&opts.Id, "id", "i", "", "email-id, only if integration type is email.")

	return cmd
}

func UpdateIntegrationRun(opts *UpdateIntegrationOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.TeamId == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
		os.Exit(0)
	}

	err = APICalls.UpdateIntegration(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.IntegrationId,
		opts.Name, opts.Description, opts.Id)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Integration updated successfully!\n", cs.SuccessIcon())
	}
}
