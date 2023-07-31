package source_delete

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"io"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type SourceDeleteOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	TeamId   string
	SourceId string
}

func NewSourceDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SourceDeleteOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Args:  cobra.ExactArgs(0),
		Short: "Delete source",
		Long: heredoc.Docf(`
			Delete a source based on the source id of a particular team.

			The user is prompted with the teams. User can select a team to show the sources.
		`),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sources delete
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive {
				if opts.TeamId == "" {
					fmt.Fprint(opts.IO.ErrOut, "team-id is required.\n")
					return
				}

				if opts.SourceId == "" {
					fmt.Fprint(opts.IO.ErrOut, "source-id is required.\n")
					return
				}
			}

			sourceDeleteRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-id", "", "Team ID for which the source is to be deleted.")
	cmd.Flags().StringVar(&opts.SourceId, "source-id", "", "Source ID for which the source is to be deleted.")
	return cmd
}

func sourceDeleteRun(opts *SourceDeleteOptions) {
	cs := opts.IO.ColorScheme()

	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config.\n", cs.FailureIcon())
		return
	}

	if opts.Interactive && opts.TeamId == "" && opts.SourceId == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		opts.SourceId, _ = pre_defined_prompters.AskSourceId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

	} else {
		if opts.TeamId == "" || opts.SourceId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s team-id and source-id both are required.\n", cs.FailureIcon())
		}
	}

	err = deleteSource(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.SourceId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s.\n", cs.FailureIcon(), err.Error())
	}

	fmt.Fprintf(opts.IO.Out, "%s Source deleted successfully.\n", cs.SuccessIcon())
}

func deleteSource(client *http.Client, token, endpoint string, teamId, sourceId string) error {
	url := fmt.Sprintf(endpoint+"api/team/%s/source/%s", teamId, sourceId)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", "Logfire-cli")

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var sourceResp models.SourceResponse

	err = json.Unmarshal(body, &sourceResp)
	if err != nil {
		return err
	}

	if !sourceResp.IsSuccessful {
		return errors.New("api error")
	}

	return nil
}
