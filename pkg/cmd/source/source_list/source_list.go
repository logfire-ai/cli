package source_list

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/source/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type SourceListOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	TeamId string
}

func NewSourceListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SourceListOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Args:  cobra.ExactArgs(0),
		Short: "Get sources",
		Long: heredoc.Docf(`
			Get sources for a particular team.

			The user is prompted with the teams. User can select a team to show the sources.
		`),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sources list
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if opts.TeamId == "" && !opts.Interactive {
				fmt.Fprint(opts.IO.ErrOut, "team-id is required.\n")
			}

			sourceListRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-id", "", "Team ID for which the sources will be fetched.")
	return cmd
}

func sourceListRun(opts *SourceListOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.TeamId == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s team-id is required.\n", cs.FailureIcon())
	}
	sources, err := getAllSources(opts.HttpClient(), cfg.Get().Token, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		return
	}

	fmt.Fprintf(opts.IO.Out, "%s Successfully fetched sources for team-id %s\n", cs.SuccessIcon(), opts.TeamId)

	for _, v := range sources {
		fmt.Fprintf(opts.IO.Out, "%s %s %s %s\n", cs.IntermediateIcon(), v.Name, v.ID, v.Platform)
	}
}

func getAllSources(client *http.Client, token, teamId string) ([]models.Source, error) {
	url := fmt.Sprintf("https://api.logfire.sh/api/team/%s/source", teamId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []models.Source{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return []models.Source{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []models.Source{}, err
	}

	var sourceResp models.SourceResponse

	err = json.Unmarshal(body, &sourceResp)
	if err != nil {
		return []models.Source{}, err
	}

	if !sourceResp.IsSuccessful {
		return []models.Source{}, errors.New("Api error!")
	}

	return sourceResp.Data, nil
}
