package bootstrap

import (
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"os"
)

type PromptBootstrapOptions struct {
	IO *iostreams.IOStreams

	Interactive bool
}

func NewCmdBootstrap(f *cmdutil.Factory) *cobra.Command {
	opts := &PromptBootstrapOptions{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:     "bootstrap",
		Short:   "bootstrap",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PromptBootstrapRun(opts)
		},
	}

	return cmd
}

func PromptBootstrapRun(opts *PromptBootstrapOptions) {
	if !opts.Interactive {
		os.Exit(1)
	}

	if opts.Interactive {
		prompter.NewOnboardingForm()
	}
}
