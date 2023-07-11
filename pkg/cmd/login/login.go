package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
		Short: "Login to logfire.sh",
		Long: heredoc.Docf(`
			Login to logfire.sh using a password or token.

			There are two ways to login to logfire.sh, using a password or by using the token
			provided in the magic link. By default the cli will give a prompt to select from.

			Alternatively, use %[1]s--with-token%[1]s to pass in a token on standard input.
			or %[1]s--with-password%[1]s to pass in a password on standard input.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire login

			# authenticate against logfire.sh by reading the password from the prompt
			$ logfire login --email name@example.com --password asdf@1234

			# authenticate against logfire.sh by reading the token from the prompt
			$ logfire login --email name@example.com --token myToken
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if opts.Password != "" && opts.Token != "" {
				fmt.Fprint(opts.IO.ErrOut, "Please provide either password or token.\n")
			}

			if !opts.Interactive {
				if opts.Email == "" {
					fmt.Fprint(opts.IO.ErrOut, "Email is required\n")
				}

				if opts.Token == "" || opts.Password == "" {
					fmt.Fprint(opts.IO.ErrOut, "Please provide either password or token for sign in.")
				}
			}

			loginRun(opts)
		},
		GroupID: "core",
	}

	cmd.Flags().StringVar(&opts.Email, "email", "", "Email ID of the user.")
	cmd.Flags().StringVar(&opts.Password, "password", "", "Password of the user.")
	cmd.Flags().StringVar(&opts.Token, "token", "", "Single Sign in token of the user.")
	return cmd
}

type SigninRequest struct {
	Email      string `json:"email,omitempty"`
	AuthType   int    `json:"authType"`
	Credential string `json:"credential"`
}

func loginRun(opts *LoginOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Email == "" {
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
	}

	opts.IO.StartProgressIndicatorWithLabel("Logging in to logfire.sh")

	resp, err := PasswordSignin(opts.Email, opts.Password)
	opts.IO.StopProgressIndicator()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s Signin Failed. %s\n", cs.FailureIcon(), err.Error())
		return
	}

	cfg.UpdateConfig(opts.Email, resp.BearerToken.AccessToken, resp.UserBody.ProfileID)
	fmt.Fprintf(opts.IO.Out, "\n%s Logged in as %s\n", cs.SuccessIcon(), cs.Bold(opts.Email))
}

func PasswordSignin(email, password string) (Response, error) {
	var response Response
	signinReq := SigninRequest{
		Email:      email,
		AuthType:   2,
		Credential: password,
	}

	reqBody, err := json.Marshal(signinReq)
	if err != nil {
		return response, err
	}

	url := "https://api.logfire.sh/api/auth/signin"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, errors.New("wrong password")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func TokenSignin(token string) (Response, error) {
	var response Response

	signinReq := SigninRequest{
		AuthType:   1,
		Credential: strings.TrimSpace(token),
	}

	reqBody, err := json.Marshal(signinReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return response, err
	}

	url := "https://api.logfire.sh/api/auth/signin"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Sign-in successful!")
	} else {
		fmt.Printf("Sign-in failed with status code: %d\n", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		return response, err
	}

	return response, nil
}
