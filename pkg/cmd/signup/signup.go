package signup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/login"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type SignupOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
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
		Short: "Signup to logfire.sh",
		Long: heredoc.Docf(`
			Signup to logfire.sh to create a new account.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire auth signup
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}
			signupRun(opts)
		},
		GroupID: "core",
	}

	return cmd
}

func signupRun(opts *SignupOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	email, err := opts.Prompter.Input("Enter your email:", "")
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read email\n", cs.FailureIcon())
		return
	}

	err = SignupFlow(email)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s Error while signing up %s\n", cs.FailureIcon(), err.Error())
		return
	}

	fmt.Fprintf(opts.IO.ErrOut, "%s Thank You for Registering. An email has been sent to your adress %s\n", cs.SuccessIcon(), cs.Bold(email))

	credentialToken, err := opts.Prompter.Input("Please paste the token in the email link here:", "")
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Unable to read token %s\n", cs.SuccessIcon(), cs.Bold(credentialToken))
		return
	}

	resp, err := login.TokenSignin(credentialToken)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s Error while signing up with the token %s", cs.FailureIcon(), err.Error())
		return
	}

	cfg.UpdateConfig(email, resp.BearerToken.AccessToken, resp.UserBody.ProfileID)

	err = OnboardingFlow(opts.IO, opts.Prompter, resp.UserBody.ProfileID, cfg.Get().Token)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "\n%s %s", cs.FailureIcon(), err.Error())
		return
	}

	fmt.Fprintf(opts.IO.Out, "%s User onboarded successfully.\n", cs.SuccessIcon())
}

func SignupFlow(email string) error {
	signupReq := SignupRequest{
		Email: email,
	}

	reqBody, err := json.Marshal(signupReq)
	if err != nil {
		return err
	}

	url := "https://api.logfire.sh/api/auth/signup"

	transport := http.Transport{
		IdleConnTimeout:   30 * time.Second,
		MaxIdleConns:      100,
		MaxConnsPerHost:   0,
		DisableKeepAlives: false,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return err
	}

	return nil
}

func OnboardingFlow(IO *iostreams.IOStreams, prompt prompter.Prompter, profileID, authToken string) error {
	cs := IO.ColorScheme()

	firstName, err := prompt.Input("Enter your first name:", "")
	if err != nil {
		return err
	}

	lastName, err := prompt.Input("Enter your last name:", "")
	if err != nil {
		return err
	}

	onboardReq := OnboardRequest{
		FirstName: firstName,
		LastName:  lastName,
	}

	reqBody, err := json.Marshal(onboardReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return err
	}

	url := "https://api.logfire.sh/api/profile/" + profileID + "/onboard"

	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("unable to process the request")
	}

	var password string
	var isConfirmed bool

	for i := 0; i < 3; i++ {
		password, err = prompt.Password("Enter the password you want to set:")
		if err != nil {
			return err
		}

		confirmPassword, err := prompt.Password("Confirm password:")
		if err != nil {
			return err
		}

		if password == confirmPassword {
			isConfirmed = true
			break
		}

		fmt.Fprintf(IO.ErrOut, "%s passwords do not match. please try again\n", cs.FailureIcon())
	}

	if !isConfirmed {
		return errors.New("maximum number of attempts exceeded")
	}

	pwdReq := SetPassword{
		Password: password,
	}

	reqBody, err = json.Marshal(pwdReq)
	if err != nil {
		return err
	}

	url = "https://api.logfire.sh/api/profile/" + profileID + "/set-password"

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	resp, err = client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("unable to set password, please try again later")
	}
	return nil
}
