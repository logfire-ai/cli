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
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
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
		},
	}

	cmd.AddCommand(livetail.NewLivetailCmd(f))
	cmd.AddCommand(view.NewViewStreamOptionsCmd(f))
	return cmd
}

func PromptStreamRun(opts *PromptStreamOptions) {
	if !opts.Interactive {
		return
	}

	if opts.Interactive {
		cs := opts.IO.ColorScheme()
		cfg, _ := config.NewConfig()

		teamid, _ := pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		err := cfg.UpdateConfig(nil, nil, nil, nil,
			nil, &teamid, nil, nil, nil)
		if err != nil {
			return
		}

		ui := gui.NewUI(cfg)
		if err := ui.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
