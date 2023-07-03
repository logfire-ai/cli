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
	"github.com/logfire-sh/cli/pkg/cmd/auth/models"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type LoginOptions struct {
	IO       *iostreams.IOStreams
	Config   config.Config
	Prompter prompter.Prompter

	Interactive bool
}

func NewLoginCmd() *cobra.Command {
	io := iostreams.System()
	cfg, _ := config.NewConfig()
	prmpt := prompter.New(io.In, io.Out, io.ErrOut)

	opts := &LoginOptions{
		IO:       io,
		Config:   *cfg,
		Prompter: prmpt,
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
			$ logfire auth login

			# authenticate against logfire.sh by reading the token from a file
			$ logfire auth login --with-token mytoken

			# authenticate with a specific GitHub instance
			$ logfire auth login --with-password your_password
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}
			loginRun(opts)
		},
	}

	return cmd
}

type SigninRequest struct {
	Email      string `json:"email,omitempty"`
	AuthType   int    `json:"authType"`
	Credential string `json:"credential"`
}

func loginRun(opts *LoginOptions) {
	cs := opts.IO.ColorScheme()

	email, err := opts.Prompter.Input("Enter your email:", "")
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read email\n", cs.FailureIcon())
		return
	}

	password, err := opts.Prompter.Password("Enter your password:")
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read password\n", cs.FailureIcon())
		return
	}

	opts.IO.StartProgressIndicatorWithLabel("Logging in to logfire.sh")

	resp, err := PasswordSignin(email, password)
	opts.IO.StopProgressIndicator()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s Signin Failed. %s\n", cs.FailureIcon(), err.Error())
		return
	}

	opts.Config.UpdateConfig(email, resp.BearerToken.AccessToken)
	fmt.Fprintf(opts.IO.Out, "\n%s Logged in as %s\n", cs.SuccessIcon(), cs.Bold(email))
}

func PasswordSignin(email, password string) (models.Response, error) {
	var response models.Response
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

func TokenSignin(token string) (models.Response, error) {
	var response models.Response

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
