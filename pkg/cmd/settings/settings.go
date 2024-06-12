package settings

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/helpers"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

var choices = []string{"Change default team", "Change theme", "Exit"}

type SettingsOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	Theme  string
	TeamId string
	Choice string
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
	client := http.Client{}

	if opts.TeamId != "" {
		opts.TeamId = helpers.TeamNameToTeamId(&client, cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		err = APICalls.UpdateFlag(cfg, cfg.Get().ProfileID, opts.TeamId, cfg.Get().EndPoint)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to update default team\n", cs.FailureIcon())
			return
		}
	}

	if opts.Interactive && opts.Theme == "" && opts.TeamId == "" {
		opts.Choice, err = opts.Prompter.Select("What do you want to do?", "", choices)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read choice\n", cs.FailureIcon())
			return
		}
	}

	if opts.Interactive && opts.Choice == "Exit" {
		os.Exit(0)
	}

	if opts.Interactive && opts.Choice == "Change theme" {
		themeOptions := []string{"Dark", "Light"}

		opts.Theme, err = opts.Prompter.Select("Select a theme:", "", themeOptions)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Theme\n", cs.FailureIcon())
			return
		}
	}

	if opts.Interactive && opts.Choice == "Change default team" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		err = APICalls.UpdateFlag(cfg, cfg.Get().ProfileID, opts.TeamId, cfg.Get().EndPoint)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to update default team\n", cs.FailureIcon())
			return
		}
	}

	if opts.Theme == "" && opts.TeamId == "" {
		fmt.Fprint(opts.IO.ErrOut, "atleast one flag is required.\n")
		return
	}

	loweredTheme := strings.ToLower(opts.Theme)

	err = cfg.UpdateConfig(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &loweredTheme)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to apply settings\n", cs.FailureIcon())
		return
	}

	fmt.Fprintf(opts.IO.Out, "\n%s Settings applied successfully\n", cs.SuccessIcon())
}
