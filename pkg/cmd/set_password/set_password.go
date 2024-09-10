package set_password

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
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type SetPasswordOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	Password string
}

func NewSetPasswordCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SetPasswordOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "set-password",
		Args:  cobra.ExactArgs(0),
		Short: "set-password for your account",
		Long: heredoc.Docf(`
			set-password for your account
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire set-password

			# set password by reading the password from the prompt
			$ logfire set-password --password <password>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			SetPasswordRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Password, "password", "p", "", "Password.")
	return cmd
}

func SetPasswordRun(opts *SetPasswordOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive && opts.Password == "" {
		opts.Password, err = opts.Prompter.Password("Enter a password:")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read password\n", cs.FailureIcon())
			return
		}
	} else {
		if opts.Password == "" {
			fmt.Fprint(opts.IO.ErrOut, "password is required.")
			os.Exit(0)
		}
	}

	err = APICalls.SetPassword(cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().ProfileID, opts.Password)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Fprintf(opts.IO.ErrOut, "%s Error: Connection failed (Server down or no internet)\n", cs.FailureIcon())
			os.Exit(0)
			return
		}
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Password set successfully!\n", cs.SuccessIcon())
	}
}
