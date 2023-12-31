package logout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		Short: "Logout from logfire.ai.",
		Long: heredoc.Docf(`
			Logout from your current account of logfire.ai.
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
		return
	}

	// logout steps

	user := cfg.Get().Username

	err = logout(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Fprintf(opts.IO.ErrOut, "%s Error: Connection failed (Server down or no internet)\n", cs.FailureIcon())
			os.Exit(0)
			return
		}
	}

	err = cfg.DeleteConfig()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to logout.\n", cs.FailureIcon())
		return
	}

	fmt.Fprintf(opts.IO.Out, "%s User %s successfully logged out.\n", cs.SuccessIcon(), user)
}

func logout(client *http.Client, token string, endpoint string, refreshToken string) error {
	url := endpoint + "api/auth/signout"

	reqBody := map[string]string{
		"AccessToken":  token,
		"RefreshToken": refreshToken,
	}

	jsonValue, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Logfire-cli")

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(body), "false") {
		// return errors.New("request can't be completed")
		return nil
	}

	return nil
}
