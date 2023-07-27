package integrations

import (
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/integrations/integrations_create"
	"github.com/logfire-sh/cli/pkg/cmd/integrations/integrations_delete"
	"github.com/logfire-sh/cli/pkg/cmd/integrations/integrations_list"
	"github.com/logfire-sh/cli/pkg/cmd/integrations/integrations_update"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

type PromptIntegrationsOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	Choice      string
}

var choices = []string{"Create", "List", "Delete", "Update"}

func NewCmdIntegrations(f *cmdutil.Factory) *cobra.Command {
	opts := &PromptIntegrationsOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:     "integrations <command>",
		Short:   "integrations",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PromptIntegrationsRun(opts)

			switch opts.Choice {
			case choices[0]:
				integrations_create.NewCreateIntegrationsCmd(f).Run(cmd, []string{})
			case choices[1]:
				integrations_list.NewListIntegrationsCmd(f).Run(cmd, []string{})
			case choices[2]:
				integrations_delete.NewDeleteIntegrationCmd(f).Run(cmd, []string{})
			case choices[3]:
				integrations_update.NewUpdateIntegrationsCmd(f).Run(cmd, []string{})
			}
		},
	}

	cmd.AddCommand(integrations_create.NewCreateIntegrationsCmd(f))
	cmd.AddCommand(integrations_list.NewListIntegrationsCmd(f))
	cmd.AddCommand(integrations_delete.NewDeleteIntegrationCmd(f))
	cmd.AddCommand(integrations_update.NewUpdateIntegrationsCmd(f))
	return cmd
}

func PromptIntegrationsRun(opts *PromptIntegrationsOptions) {
	cs := opts.IO.ColorScheme()
	if !opts.Interactive {
		return
	}

	if opts.Interactive {
		err := errors.New("")
		opts.Choice, err = opts.Prompter.Select("What do you want to do?", "", choices)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read choice\n", cs.FailureIcon())
			return
		}
	}
}
