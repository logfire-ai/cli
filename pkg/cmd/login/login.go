package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/logfire-sh/cli/pkg/cmd/login/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"

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
	Staging  bool
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
	cmd.Flags().BoolVarP(&opts.Staging, "staging", "s", false, "Use staging server?")
	return cmd
}

func loginRun(opts *LoginOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Staging {
		endpoint := "https://api-stg.logfire.ai/"
		grpc_endpoint := "api-stg.logfire.ai:443"
		grpc_ingestion := "https://in-stg.logfire.ai"

		err = cfg.UpdateConfig(nil, nil, nil, nil,
			nil, nil, &endpoint, &grpc_endpoint, &grpc_ingestion)
		if err != nil {
			return
		}
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
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read token\n", cs.FailureIcon())
				return
			}

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

	client := &http.Client{}

	reqBody, err := json.Marshal(signinReq)
	if err != nil {
		return
	}

	url := endpoint + "api/auth/signin"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	if !response.IsSuccessful {
		fmt.Fprintf(io.ErrOut, "\n%s %s\n", cs.FailureIcon(), response.Message[0])
		os.Exit(0)
	}

	io.StopProgressIndicator()

	err = cfg.UpdateConfig(&response.UserBody.Email, &response.UserBody.Role, &response.BearerToken.AccessToken, &response.UserBody.ProfileID,
		&response.BearerToken.RefreshToken, &response.UserBody.TeamID, nil, nil, nil)
	if err != nil {
		return
	}
	fmt.Fprintf(io.Out, "\n%s Logged in as %s\n", cs.SuccessIcon(), cs.Bold(response.UserBody.Email))
}

func TokenSignin(IO *iostreams.IOStreams, cfg config.Config, cs *iostreams.ColorScheme, token, endpoint string) error {
	var response models.Response

	signinReq := models.SigninRequest{
		AuthType:   1,
		Credential: strings.TrimSpace(token),
	}

	client := &http.Client{}

	reqBody, err := json.Marshal(signinReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return err
	}

	url := endpoint + "api/auth/signin"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

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
		if ok, _ := regexp.MatchString("invalid UUID", response.Message[0]); ok {
			fmt.Fprintf(IO.ErrOut, "\n%s %s\n", cs.FailureIcon(), "invalid token")
		} else {
			fmt.Fprintf(IO.ErrOut, "\n%s %s\n", cs.FailureIcon(), response.Message[0])
		}
		return errors.New(response.Message[0])
	}

	err = cfg.UpdateConfig(&response.UserBody.Email, &response.UserBody.Role, &response.BearerToken.AccessToken, &response.UserBody.ProfileID,
		&response.BearerToken.RefreshToken, &response.UserBody.TeamID, nil, nil, nil)
	if err != nil {
		return err
	}

	fmt.Fprintf(IO.Out, "\n%s Logged in as %s\n", cs.SuccessIcon(), cs.Bold(response.Email))

	return nil
}
