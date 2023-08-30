package signup

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/login"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
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
	Role            string
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
			SignupRun(opts)
		},
		GroupID: "core",
	}

	cmd.Flags().StringVarP(&opts.Email, "email", "e", "", "Email address")
	cmd.Flags().StringVarP(&opts.credentialToken, "token", "t", "", "Token received on email")
	cmd.Flags().StringVarP(&opts.FirstName, "first-name", "f", "", "First name")
	cmd.Flags().StringVarP(&opts.LastName, "last-name", "l", "", "Last name")
	cmd.Flags().StringVarP(&opts.Role, "role", "r", "", "Role")

	return cmd
}

func SignupRun(opts *SignupOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	var email string
	var interactive bool

	isEmpty := func(s string) bool {
		return s == ""
	}

	if !opts.Interactive && isEmpty(opts.Email) {
		fmt.Fprintf(opts.IO.ErrOut, "%s Email address is required\n", cs.FailureIcon())
		return
	}

	if opts.Interactive && isEmpty(opts.Email) {
		interactive = true
		opts.Email, err = opts.Prompter.Input("Enter your email:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read email\n", cs.FailureIcon())
			return
		}
	}

	if isEmpty(opts.Email) && !isEmpty(opts.credentialToken) && !interactive {
		err = OnboardingRun(opts)
		if err != nil {
			return
		}

		successMessage := fmt.Sprintf("%s User onboarded successfully.\n", cs.SuccessIcon())
		fmt.Fprint(opts.IO.Out, successMessage)
		passwordMessage := fmt.Sprintf("%s You can set your password using %s anytime later.\n", cs.SuccessIcon(), cs.Blue("\"logfire set-password --password <password>\""))
		fmt.Fprint(opts.IO.Out, passwordMessage)

		os.Exit(0)
	}

	msg, err := APICalls.SignupFlow(opts.Email, cfg.Get().EndPoint)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s Error while signing up %s\n", cs.FailureIcon(), err.Error())
		return
	}

	if msg == "already registered user. Sent link to login" {
		fmt.Printf("%s You are already a user, Sent a link to login\n", cs.SuccessIcon())

		if opts.credentialToken == "" {
			opts.credentialToken, err = opts.Prompter.Input("Please paste the token received in the email here:", "")
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Unable to read token %s\n", cs.SuccessIcon(), cs.Bold(opts.credentialToken))
			}
		}

		err = login.TokenSignin(opts.IO, cfg, cs, opts.credentialToken, cfg.Get().EndPoint)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Unable to sign in with token \n", cs.FailureIcon())
			return
		}

		return
	} else {
		registerMessage := fmt.Sprintf("%s Thank You for Registering. An email has been sent to your address %s\n", cs.SuccessIcon(), cs.Bold(email))
		fmt.Fprint(opts.IO.ErrOut, registerMessage)
	}

	if !interactive {
		onboardingMessage := fmt.Sprintf("Use %s to complete your onboarding process \n", cs.Blue("\"logfire signup --token <token received on email> --first-name <first-name> --last-name <last-name> --role <role>\""))
		fmt.Fprint(opts.IO.ErrOut, onboardingMessage)
	}

	if interactive {
		err = OnboardingRun(opts)
		if err != nil {
			return
		}

		successMessage := fmt.Sprintf("%s User onboarded successfully.\n", cs.SuccessIcon())
		fmt.Fprint(opts.IO.Out, successMessage)
		//passwordMessage := fmt.Sprintf("%s You can set your password using %s anytime later.\n", cs.SuccessIcon(), cs.Blue("\"logfire set-password --password <password>\""))
		//fmt.Fprint(opts.IO.Out, passwordMessage)
	}
}

func OnboardingRun(opts *SignupOptions) error {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Config failed to load %s\n", cs.FailureIcon(), cs.Bold(opts.credentialToken))
		return err
	}

	if opts.credentialToken == "" {
		opts.credentialToken, err = opts.Prompter.Input("Please paste the token received in the email here:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Unable to read token %s\n", cs.FailureIcon(), cs.Bold(opts.credentialToken))
			return err
		}
	}

	err = login.TokenSignin(opts.IO, cfg, cs, opts.credentialToken, cfg.Get().EndPoint)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Unable to sign in with token \n", cs.FailureIcon())
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

		if opts.Role == "" {
			opts.Role, err = opts.Prompter.Input("Enter your Role:", "")
			if err != nil {
				return err
			}
		}
	}

	err = APICalls.OnboardingFlow(cfg.Get().ProfileID, cfg.Get().Token, cfg.Get().EndPoint, opts.FirstName, opts.LastName, opts.Role)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s %s", cs.FailureIcon(), err.Error())
		return err
	}

	return nil
}
