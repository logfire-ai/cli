package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmd/login/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type LoginOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	Email    string
	Password string
	Token    string
}

func NewLoginCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LoginOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "login",
		Args:  cobra.ExactArgs(0),
		Short: "Login to logfire.ai",
		Long: heredoc.Docf(`
			Login to logfire.ai using a email and password or token.

			There are two ways to login to logfire.ai, using a password or by using the token
			provided in the magic link. By default the cli will give a prompt to select from.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire login

			# authenticate against logfire.ai by email and password
			$ logfire login --email <name@example.com> --password <password>

			# authenticate against logfire.ai by magic link token
				# First request a Magic link to your email address
				$ logfire login --email <name@example.com>
	
				# Second authenticate using the token received on your email address
				$ logfire login --token <token>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			loginRun(opts)
		},
		GroupID: "core",
	}

	cmd.Flags().StringVarP(&opts.Email, "email", "e", "", "Email ID of the user.")
	cmd.Flags().StringVarP(&opts.Password, "password", "p", "", "Password of the user.")
	cmd.Flags().StringVarP(&opts.Token, "token", "t", "", "Single Sign in token of the user.")
	return cmd
}

func loginRun(opts *LoginOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	var choiceList = []string{"Magic link", "Password"}

	if opts.Interactive && opts.Token == "" && opts.Email == "" && opts.Password == "" {

		choice, err := opts.Prompter.Select("Select login method (Default: Magic link)", "Magic link", choiceList)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read choice\n", cs.FailureIcon())
			return
		}

		switch choice {
		case "Magic link":
			opts.Email, err = opts.Prompter.Input("Enter your email:", "")
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read email\n", cs.FailureIcon())
				return
			}

			err = APICalls.SendMagicLink(cfg.Get().EndPoint, opts.Email)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to send magic link\n", cs.FailureIcon())
				return
			}

			fmt.Fprintf(opts.IO.Out, "%s Magic link sent to %s\n", cs.SuccessIcon(), opts.Email)

			opts.Token, err = opts.Prompter.Input("Enter the token received on your email:", "")

			opts.IO.StartProgressIndicatorWithLabel("Logging in to logfire.ai")
			TokenSignin(opts.IO, cfg, cs, opts.Token, cfg.Get().EndPoint)

		case "Password":
			email, err := opts.Prompter.Input("Enter your email:", "")
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read email\n", cs.FailureIcon())
				return
			}
			opts.Email = email

			password, err := opts.Prompter.Password("Enter your password:")
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read password\n", cs.FailureIcon())
				return
			}
			opts.Password = password

			opts.IO.StartProgressIndicatorWithLabel("Logging in to logfire.ai")
			PasswordSignin(opts.IO, cfg, cs, opts.Email, opts.Password, cfg.Get().EndPoint)
		}

	} else {
		isEmpty := func(s string) bool {
			return s == ""
		}

		switch {
		case !isEmpty(opts.Token) && isEmpty(opts.Email) && isEmpty(opts.Password):
			opts.IO.StartProgressIndicatorWithLabel("Logging in to logfire.ai")
			TokenSignin(opts.IO, cfg, cs, opts.Token, cfg.Get().EndPoint)

		case !isEmpty(opts.Email) && isEmpty(opts.Token):
			if isEmpty(opts.Password) {
				err = APICalls.SendMagicLink(cfg.Get().EndPoint, opts.Email)
				if err != nil {
					fmt.Fprintf(opts.IO.ErrOut, "%s Failed to send magic link\n", cs.FailureIcon())
					return
				}
				fmt.Fprintf(opts.IO.Out, "%s Magic link sent to %s \n%s Use \"logfire login --token <token>\" to sign-in using received token\n", cs.SuccessIcon(), opts.Email, cs.IntermediateIcon())
			} else {
				opts.IO.StartProgressIndicatorWithLabel("Logging in to logfire.ai")
				PasswordSignin(opts.IO, cfg, cs, opts.Email, opts.Password, cfg.Get().EndPoint)
			}

		case !isEmpty(opts.Email) && !isEmpty(opts.Password) && !isEmpty(opts.Token):
			fmt.Fprint(opts.IO.ErrOut, "Please provide only one method of authentication: either email and password or token\n")

		case isEmpty(opts.Token) && isEmpty(opts.Email) && isEmpty(opts.Password):
			fmt.Fprint(opts.IO.ErrOut, "Please provide either email and password or token\n")

		case !isEmpty(opts.Email) && !isEmpty(opts.Token) && isEmpty(opts.Password):
			fmt.Fprint(opts.IO.ErrOut, "Please provide either email and password or token, not both\n")

		case !isEmpty(opts.Password) && !isEmpty(opts.Token) && isEmpty(opts.Email):
			fmt.Fprint(opts.IO.ErrOut, "Please provide either email and password or token, not both\n")
		}
	}
}

func PasswordSignin(io *iostreams.IOStreams, cfg config.Config, cs *iostreams.ColorScheme, email, password string, endpoint string) {
	var response models.Response
	signinReq := models.SigninRequest{
		Email:      email,
		AuthType:   2,
		Credential: password,
	}

	reqBody, err := json.Marshal(signinReq)
	if err != nil {
		return
	}

	url := endpoint + "api/auth/signin"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	if !response.IsSuccessful {
		fmt.Fprintf(io.ErrOut, "\n%s %s %s\n", cs.FailureIcon(), response.Message[0], err.Error())
	}

	io.StopProgressIndicator()

	cfg.UpdateConfig(&response.UserBody.Email, &response.BearerToken.AccessToken, &response.UserBody.ProfileID,
		&response.BearerToken.RefreshToken, nil)
	fmt.Fprintf(io.Out, "\n%s Logged in as %s\n", cs.SuccessIcon(), cs.Bold(response.UserBody.Email))

	return
}

func TokenSignin(IO *iostreams.IOStreams, cfg config.Config, cs *iostreams.ColorScheme, token, endpoint string) error {
	var response models.Response

	signinReq := models.SigninRequest{
		AuthType:   1,
		Credential: strings.TrimSpace(token),
	}

	reqBody, err := json.Marshal(signinReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return err
	}

	url := endpoint + "api/auth/signin"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	IO.StopProgressIndicator()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		return err
	}

	if !response.IsSuccessful {
		fmt.Fprintf(IO.ErrOut, "\n%s %s\n", cs.FailureIcon(), response.Message[0])
		return errors.New(response.Message[0])
	}

	err = cfg.UpdateConfig(&response.UserBody.Email, &response.BearerToken.AccessToken, &response.UserBody.ProfileID,
		&response.BearerToken.RefreshToken, nil)
	if err != nil {
		return err
	}

	fmt.Fprintf(IO.Out, "\n%s Logged in as %s\n", cs.SuccessIcon(), cs.Bold(response.Email))

	return nil
}
