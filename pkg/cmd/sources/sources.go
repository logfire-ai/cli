package sources

import (
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/sources/source_create"
	"github.com/logfire-sh/cli/pkg/cmd/sources/source_delete"
	"github.com/logfire-sh/cli/pkg/cmd/sources/source_list"
	"github.com/logfire-sh/cli/pkg/cmd/sources/source_update"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

type PromptSourceOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	Choice      string
}

var choices = []string{"Create", "List", "Delete", "Update"}

func NewCmdSource(f *cmdutil.Factory) *cobra.Command {
	opts := &PromptSourceOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:     "sources <command>",
		Short:   "Get source for a team",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PromptSourceRun(opts)

			switch opts.Choice {
			case choices[0]:
				source_create.NewSourceCreateCmd(f).Run(cmd, []string{})
			case choices[1]:
				source_list.NewSourceListCmd(f).Run(cmd, []string{})
			case choices[2]:
				source_delete.NewSourceDeleteCmd(f).Run(cmd, []string{})
			case choices[3]:
				source_update.NewSourceUpdateCmd(f).Run(cmd, []string{})
			}
		},
	}

	cmd.AddCommand(source_list.NewSourceListCmd(f))
	cmd.AddCommand(source_create.NewSourceCreateCmd(f))
	cmd.AddCommand(source_update.NewSourceUpdateCmd(f))
	cmd.AddCommand(source_delete.NewSourceDeleteCmd(f))
	return cmd
}

func PromptSourceRun(opts *PromptSourceOptions) {
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
