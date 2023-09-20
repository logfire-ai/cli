package settings

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

type SettingsOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	Theme string
}

func SettingsCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SettingsOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "settings",
		Args:  cobra.ExactArgs(0),
		Short: "settings",
		Long: heredoc.Docf(`
			Change CLI Settings.
		`),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire settings

			# start argument setup
			$ logfire settings --theme <light | dark>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			SettingsRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Theme, "theme", "t", "", "Light Or Dark theme (Only for 'logfire stream' command)")
	return cmd
}

func SettingsRun(opts *SettingsOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Interactive && opts.Theme == "" {
		themeOptions := []string{"Dark", "Light"}

		opts.Theme, err = opts.Prompter.Select("Select a theme:", "", themeOptions)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Theme\n", cs.FailureIcon())
			return
		}
	} else {
		if opts.Theme == "" {
			fmt.Fprint(opts.IO.ErrOut, "theme is required.\n")
			return
		}
	}

	loweredTheme := strings.ToLower(opts.Theme)

	err = cfg.UpdateConfig(nil, nil, nil, nil, nil, nil, nil, nil, nil, &loweredTheme)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to apply settings\n", cs.FailureIcon())
		return
	}

	fmt.Fprintf(opts.IO.Out, "\n%s Settings applied successfully\n", cs.SuccessIcon())
}
