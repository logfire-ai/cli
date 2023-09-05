package stream

import (
	"log"
	"net/http"

	"github.com/logfire-sh/cli/gui"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/stream/livetail"
	"github.com/logfire-sh/cli/pkg/cmd/stream/view"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type PromptStreamOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	Choice      string
}

func NewCmdStream(f *cmdutil.Factory, cmdCh chan bool) *cobra.Command {
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

			PromptStreamRun(opts, cmdCh)
		},
	}

	cmd.AddCommand(livetail.NewLivetailCmd(f))
	cmd.AddCommand(view.NewViewStreamOptionsCmd(f))
	return cmd
}

func PromptStreamRun(opts *PromptStreamOptions, cmdCh chan bool) {
	if !opts.Interactive {
		return
	}

	if opts.Interactive {
		ui := gui.NewUI(cmdCh)
		if err := ui.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
