package delete_profile

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type DeleteProfileOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)
}

func DeleteProfileCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteProfileOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "delete-profile",
		Args:  cobra.ExactArgs(0),
		Short: "delete your profile",
		Long: heredoc.Docf(`
			delete your profile
		`, "`"),
		Example: heredoc.Doc(`
			# delete your profile
			$ logfire delete-profile
		`),
		Run: func(cmd *cobra.Command, args []string) {
			DeleteProfileRun(opts)
		},
	}

	return cmd
}

func DeleteProfileRun(opts *DeleteProfileOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	err = APICalls.DeleteProfile(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().ProfileID)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Profile deleted successfully!\n", cs.SuccessIcon())
	}
}
