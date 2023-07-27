package views

import (
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/views/views_delete"
	"github.com/logfire-sh/cli/pkg/cmd/views/views_list"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

type PromptViewsOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	Choice      string
}

var choices = []string{"List", "Delete"}

func NewCmdViews(f *cmdutil.Factory) *cobra.Command {
	opts := &PromptViewsOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:     "views <command>",
		Short:   "Views",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PromptViewsRun(opts)

			switch opts.Choice {
			case choices[0]:
				views_list.NewViewListCmd(f).Run(cmd, []string{})
			case choices[1]:
				views_delete.NewDeleteCmd(f).Run(cmd, []string{})
			}
		},
	}

	cmd.AddCommand(views_delete.NewDeleteCmd(f))
	cmd.AddCommand(views_list.NewViewListCmd(f))
	return cmd
}

func PromptViewsRun(opts *PromptViewsOptions) {
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
