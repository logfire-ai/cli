package alerts

import (
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_create"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_delete"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_list"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_pause"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/alerts_update"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

type PromptAlertOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	Choice      string
}

var choices = []string{"create", "list", "delete", "pause", "update"}

func NewCmdAlerts(f *cmdutil.Factory) *cobra.Command {
	opts := &PromptAlertOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:     "alerts <command>",
		Short:   "alerts",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PromptAlertRun(opts)

			switch opts.Choice {
			case choices[0]:
				alerts_create.NewCreateAlertCmd(f).Run(cmd, []string{})
			case choices[1]:
				alerts_list.NewListAlertCmd(f).Run(cmd, []string{})
			case choices[2]:
				alerts_delete.NewDeleteAlertCmd(f).Run(cmd, []string{})
			case choices[3]:
				alerts_pause.NewPauseAlertCmd(f).Run(cmd, []string{})
			case choices[4]:
				alerts_update.NewAlertUpdateCmd(f).Run(cmd, []string{})
			}
		},
	}

	cmd.AddCommand(alerts_create.NewCreateAlertCmd(f))
	cmd.AddCommand(alerts_list.NewListAlertCmd(f))
	cmd.AddCommand(alerts_delete.NewDeleteAlertCmd(f))
	cmd.AddCommand(alerts_pause.NewPauseAlertCmd(f))
	cmd.AddCommand(alerts_update.NewAlertUpdateCmd(f))
	return cmd
}

func PromptAlertRun(opts *PromptAlertOptions) {
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
