package stream

import (
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/stream/livetail"
	"github.com/logfire-sh/cli/pkg/cmd/stream/view"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

type PromptStreamOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	Choice      string
}

var choices = []string{"Stream livetail", "Stream a saved view"}

func NewCmdStream(f *cmdutil.Factory) *cobra.Command {
	opts := &PromptStreamOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:     "stream",
		Short:   "Logs streaming",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PromptStreamRun(opts)

			switch opts.Choice {
			case choices[0]:
				livetail.NewLivetailCmd(f).Run(cmd, []string{})
			case choices[1]:
				view.NewViewStreamOptionsCmd(f).Run(cmd, []string{})
			}
		},
	}

	cmd.AddCommand(livetail.NewLivetailCmd(f))
	cmd.AddCommand(view.NewViewStreamOptionsCmd(f))
	return cmd
}

func PromptStreamRun(opts *PromptStreamOptions) {
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
