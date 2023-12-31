package source_create

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
	"github.com/logfire-sh/cli/pkg/cmd/source/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

var platformMap = map[string]int{
	"kubernetes": 1,
}

type SourceListOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	TeamId     string
	SourceName string
	Platform   string
}

func NewSourceCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SourceListOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Args:  cobra.ExactArgs(0),
		Short: "Create source",
		Long: heredoc.Docf(`
			Create a source for the particular team.
		`),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sources create
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

				if opts.SourceName == "" {
					fmt.Fprint(opts.IO.ErrOut, "name is required.\n")
					return
				}

				if opts.Platform == "" {
					fmt.Fprint(opts.IO.ErrOut, "platform is required.\n")
					return
				}
			}

			sourceListRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-id", "", "Team ID for which the source will be created.")
	cmd.Flags().StringVar(&opts.SourceName, "name", "", "Name of he source to be created.")
	cmd.Flags().StringVar(&opts.Platform, "platform", "", "Platform name for which you want to create source.")
	return cmd
}

func sourceListRun(opts *SourceListOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	// TODO: Add interactive flow

	if opts.TeamId == "" || opts.SourceName == "" || opts.Platform == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s team-id, name and plaform are required.\n", cs.FailureIcon())
		return
	}
	source, err := createSource(opts.HttpClient(), cfg.Get().Token, opts.TeamId, opts.SourceName, opts.Platform)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		return
	}

	fmt.Fprintf(opts.IO.Out, "%s Successfully created source for team-id %s\n", cs.SuccessIcon(), opts.TeamId)
	fmt.Fprintf(opts.IO.Out, "%s %s %s %s %s\n", cs.IntermediateIcon(), source.Name, source.ID, source.SourceToken, source.Platform)
}

func createSource(client *http.Client, token, teamId, sourceName, platform string) (models.Source, error) {
	// platform should be mapped to its respective int as sourceType, for kubernetes its 1
	sourceType, exists := platformMap[strings.ToLower(platform)]
	if !exists {
		return models.Source{}, errors.New("invalid platform")
	}

	data := models.SourceCreate{
		Name:       sourceName,
		SourceType: sourceType,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return models.Source{}, err
	}

	url := "https://api.logfire.sh/api/team/" + teamId + "/source"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.Source{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return models.Source{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Source{}, err
	}

	var sourceResp models.SourceCreateResponse
	err = json.Unmarshal(body, &sourceResp)
	if err != nil {
		return models.Source{}, err
	}

	if !sourceResp.IsSuccessful {
		fmt.Print(sourceResp)
		return models.Source{}, errors.New("failed to create source")
	}

	return sourceResp.Data, nil
}
