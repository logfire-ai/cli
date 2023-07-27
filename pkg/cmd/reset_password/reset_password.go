package reset_password

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

type ResetPasswordOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	Password string
}

func NewResetPasswordCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ResetPasswordOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "reset-password",
		Args:  cobra.ExactArgs(0),
		Short: "reset-password for your account",
		Long: heredoc.Docf(`
			reset-password for your account
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire reset-password

			# reset password by reading the password from the prompt
			$ logfire reset-password --password <password>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			ResetPasswordRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Password, "password", "", "Password of the user.")
	return cmd
}

func ResetPasswordRun(opts *ResetPasswordOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive && opts.Password == "" {
		opts.Password, err = opts.Prompter.Password("Enter a new password:")
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

	err = APICalls.ResetPassword(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().ProfileID, opts.Password)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Password reset successfully!\n", cs.SuccessIcon())
	}
}
