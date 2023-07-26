package signup

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/login"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

type SignupOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	Email           string
	credentialToken string
	FirstName       string
	LastName        string
}

func NewSignupCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SignupOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config:     f.Config,
		Prompter:   f.Prompter,
	}

	cmd := &cobra.Command{
		Use:   "signup",
		Args:  cobra.ExactArgs(0),
		Short: "Signup to logfire",
		Long: heredoc.Docf(`
			Signup to logfire to create a new account.
		`, "`"),
		Example: heredoc.Doc(`
			# start normal setup
			$ logfire signup --email <email>
			$ logfire signup --token <token received on email> --first-name <first-name> --last-name <last-name>


			# start interactive setup
			$ logfire signup
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}
			signupRun(opts)
		},
		GroupID: "core",
	}

	cmd.Flags().StringVarP(&opts.Email, "email", "e", "", "Email address")
	cmd.Flags().StringVarP(&opts.credentialToken, "token", "t", "", "Token received on email")
	cmd.Flags().StringVarP(&opts.FirstName, "first-name", "f", "", "First name")
	cmd.Flags().StringVarP(&opts.LastName, "last-name", "l", "", "Last name")

	return cmd
}

func signupRun(opts *SignupOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	var email string
	var interactive bool

	if opts.Email == "" && opts.credentialToken != "" && !interactive {
		err = OnboardingRun(opts)
		if err != nil {
			return
		}

		fmt.Fprintf(opts.IO.Out, "%s User onboarded successfully.\n", cs.SuccessIcon())
		fmt.Fprintf(opts.IO.Out, "%s You can set your password using %s anytime later.\n", cs.SuccessIcon(), cs.Blue("\"logfire set-password --password <password>\""))

		os.Exit(0)
	}

	if opts.Interactive && opts.Email == "" {
		interactive = true
		opts.Email, err = opts.Prompter.Input("Enter your email:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read email\n", cs.FailureIcon())
			return
		}
	}

	if !opts.Interactive && opts.Email == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Email address is required\n", cs.FailureIcon())
		return
	}

	err = APICalls.SignupFlow(opts.Email, cfg.Get().EndPoint)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s Error while signing up %s\n", cs.FailureIcon(), err.Error())
		return
	}

	if interactive {
		fmt.Fprintf(opts.IO.ErrOut, "%s Thank You for Registering. An email has been sent to your address %s\n", cs.SuccessIcon(), cs.Bold(email))
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Thank You for Registering. An email has been sent to your address %s\n", cs.SuccessIcon(), cs.Bold(email))
		fmt.Fprintf(opts.IO.ErrOut, "Use %s to complete your onboarding process \n", cs.Blue("\"logfire signup --token <token received on email> --first-name <first-name> --last-name <last-name>\""))
	}

	if interactive {
		err = OnboardingRun(opts)
		if err != nil {
			return
		}

		fmt.Fprintf(opts.IO.Out, "%s User onboarded successfully.\n", cs.SuccessIcon())
		fmt.Fprintf(opts.IO.Out, "%s You can set your password using %s anytime later.\n", cs.SuccessIcon(), cs.Blue("\"logfire set-password --password <password>\""))
	}
}

func OnboardingRun(opts *SignupOptions) error {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()

	if opts.credentialToken == "" {
		opts.credentialToken, err = opts.Prompter.Input("Please paste the token in the email link here:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Unable to read token %s\n", cs.SuccessIcon(), cs.Bold(opts.credentialToken))
			return err
		}
	}

	resp, err := login.TokenSignin(opts.IO, opts.credentialToken, cfg.Get().EndPoint)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s Error while signing up with the token %s", cs.FailureIcon(), err.Error())
		return err
	}

	err = cfg.UpdateConfig(resp.UserBody.Email, resp.BearerToken.AccessToken, resp.UserBody.ProfileID, resp.BearerToken.RefreshToken)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s Error updating config %s", cs.FailureIcon(), err.Error())
		return err
	}

	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return err
	}

	if opts.FirstName == "" {
		opts.FirstName, err = opts.Prompter.Input("Enter your first name:", "")
		if err != nil {
			return err
		}

		if opts.LastName == "" {
			opts.LastName, err = opts.Prompter.Input("Enter your last name:", "")
			if err != nil {
				return err
			}
		}
	}

	err = APICalls.OnboardingFlow(opts.IO, opts.Prompter, resp.UserBody.ProfileID, cfg.Get().Token, cfg.Get().EndPoint, opts.FirstName, opts.LastName)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s %s", cs.FailureIcon(), err.Error())
		return err
	}

	return nil
}
