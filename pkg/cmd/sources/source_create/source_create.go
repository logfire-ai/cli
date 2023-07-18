package source_create

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type SourceType uint8

// Declare related constants for each SourceType starting with index 0
const (
	Kubernetes SourceType = iota + 1 // EnumIndex = 1
	AWS                              // EnumIndex = 2
	JavaScript                       // EnumIndex = 3
	Docker
	Nginx
	Dokku
	FlyDotio
	Heroku
	Ubuntu
	Vercel
	DotNET
	Apache2
	Cloudflare
	Java
	Python
	PHP
	PostgreSQL
	Redis
	Ruby
	MongoDB
	MySQL
	HTTP
	Vector
	FluentBit
	Fluentd
	Logstash
	RSyslog
	Render
	SyslogNg
	Demo
)

var enumMap map[SourceType]string = map[SourceType]string{
	Kubernetes: "kubernetes",
	AWS:        "aws",
	JavaScript: "javascript",
	Docker:     "docker",
	Nginx:      "nginx",
	Dokku:      "dokku",
	FlyDotio:   "fly.io",
	Heroku:     "heroku",
	Ubuntu:     "ubuntu",
	Vercel:     "vercel",
	DotNET:     ".net",
	Apache2:    "apache2",
	Cloudflare: "cloudflare",
	Java:       "java",
	Python:     "python",
	PHP:        "php",
	PostgreSQL: "postgresql",
	Redis:      "redis",
	Ruby:       "ruby",
	MongoDB:    "mongodb",
	MySQL:      "mysql",
	HTTP:       "http",
	Vector:     "vector",
	FluentBit:  "fluentbit",
	Fluentd:    "fluentd",
	Logstash:   "logstash",
	RSyslog:    "rsyslog",
	Render:     "render",
	SyslogNg:   "syslog-ng",
	Demo:       "demo",
}

var platformMap map[string]int = map[string]int{
	"kubernetes": 1,
	"aws":        2,
	"javascript": 3,
	"docker":     4,
	"nginx":      5,
	"dokku":      6,
	"fly.io":     7,
	"heroku":     8,
	"ubuntu":     9,
	"vercel":     10,
	".net":       11,
	"apache2":    12,
	"cloudflare": 13,
	"java":       13,
	"python":     14,
	"php":        15,
	"postgresql": 16,
	"redis":      17,
	"ruby":       18,
	"mongodb":    19,
	"mysql":      20,
	"http":       21,
	"vector":     22,
	"fluentbit":  23,
	"fluentd":    24,
	"logstash":   25,
	"rsyslog":    26,
	"render":     27,
	"syslog-ng":  28,
	"demo":       29,
}

func (d SourceType) String() string {

	if d < Kubernetes || d > Demo {
		return "Unknown"
	}
	return enumMap[d]
}

func (d SourceType) EnumIndex() int {
	return int(d)
}

var platformOptions = make([]string, 0, len(platformMap))

func platformMapToArray() {
	for k := range platformMap {
		platformOptions = append(platformOptions, k)
	}
}

type SourceCreateOptions struct {
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
	opts := &SourceCreateOptions{
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

			# start argument setup
			$ logfire sources create --teamid <team-id> --name <source-name> --platform <platform>
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

			sourceCreateRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-id", "", "Team ID for which the source will be created.")
	cmd.Flags().StringVar(&opts.SourceName, "name", "", "Name of the source to be created.")
	cmd.Flags().StringVar(&opts.Platform, "platform", "", "Platform name for which you want to create source.")
	return cmd
}

func sourceCreateRun(opts *SourceCreateOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Interactive {
		opts.TeamId, err = opts.Prompter.Input("Enter TeamId:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read TeamId\n", cs.FailureIcon())
			return
		}

		opts.SourceName, err = opts.Prompter.Input("Enter Source name:", "")
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Source name\n", cs.FailureIcon())
			return
		}

		platformMapToArray()
		intPlatform, err := opts.Prompter.Select("Enter Platform name:", "", platformOptions)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Platform name\n", cs.FailureIcon())
			return
		}

		opts.Platform = platformOptions[intPlatform]
	}

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
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
