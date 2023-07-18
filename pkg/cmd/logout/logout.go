package logout

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

type LogoutOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
}

func NewLogoutCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LogoutOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "logout",
		Args:  cobra.ExactArgs(0),
		Short: "Logout from logfire.sh.",
		Long: heredoc.Docf(`
			Logout from your current account of logfire.sh.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire logout
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			logoutRun(opts)
		},
		GroupID: "core",
	}

	return cmd
}

func logoutRun(opts *LogoutOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	// Check if logged in
	if cfg.Get().ProfileID == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s No user is currently logged in.\n", cs.FailureIcon())
		return
	}

	// logout steps

	user := cfg.Get().Username

	err = logout(opts.HttpClient(), cfg.Get().Token, cfg.Get().RefreshToken)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s.\n", cs.FailureIcon(), err.Error())
		return
	}

	err = cfg.DeleteConfig()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to logout.\n", cs.FailureIcon())
		return
	}

	fmt.Fprintf(opts.IO.Out, "%s User %s successfully logged out.\n", cs.SuccessIcon(), user)
}

func logout(client *http.Client, token string, refreshToken string) error {
	url := "https://api.logfire.sh/api/auth/signout"

	reqBody := map[string]string{
		"AccessToken":  token,
		"RefreshToken": refreshToken,
	}

	jsonValue, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(body), "false") {
		return errors.New("request can't be completed")
	}

	return nil
}
