package roundtrip

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var platformOptions = []string{
	"Kubernetes",
	"AWS",
	"JavaScript",
	"Docker",
	"Nginx",
	"Dokku",
	"Fly.io",
	"Heroku",
	"Ubuntu",
	"Vercel",
	".Net",
	"Apache2",
	"Cloudflare",
	"Java",
	"Python",
	"PHP",
	"PostgreSQL",
	"Redis",
	"Ruby",
	"Mongodb",
	"MySQL",
	"HTTP",
	"Vector",
	"fluentbit",
	"Fluentd",
	"Logstash",
	"Rsyslog",
	"Render",
	"syslog-ng",
}

type PromptRoundTripOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	Choice      string

	TeamId     string
	SourceName string
	Platform   string
}

var choices = []string{"Create", "List", "Delete", "Update"}

func NewCmdRoundTrip(f *cmdutil.Factory) *cobra.Command {
	opts := &PromptRoundTripOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:     "roundtrip",
		Short:   "roundtrip",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PromptRoundTripRun(f, cmd, opts)
		},
	}

	return cmd
}

var stop = make(chan bool)

func PromptRoundTripRun(f *cmdutil.Factory, cmd *cobra.Command, opts *PromptRoundTripOptions) {
	cfg, _ := opts.Config()
	cs := opts.IO.ColorScheme()
	if !opts.Interactive {
		return
	}

	if opts.Interactive {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		sourceList, err := APICalls.GetAllSources(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err)
			return
		}

		var sourceIdNames []string

		for _, source := range sourceList {
			sourceIdNames = append(sourceIdNames, source.Name+" - "+source.ID+" - "+source.SourceToken)
		}

		var selectedSource string
		var sourceId string
		var sourceToken string

		if len(sourceList) != 0 {
			selectedSource, err = opts.Prompter.Select("Select Source to round trip:", "", sourceIdNames)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to get sources\n", cs.FailureIcon())
				return
			}

			sourceId = strings.TrimSpace(strings.Split(selectedSource, " - ")[1])
			sourceToken = strings.TrimSpace(strings.Split(selectedSource, " - ")[2])
		} else {
			fmt.Fprintf(opts.IO.ErrOut, "%s\n", "Seems that you have no sources in this team, please create a new one")

			opts.SourceName, err = opts.Prompter.Input("Enter Source name:", "")
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Source name\n", cs.FailureIcon())
				return
			}

			opts.Platform, err = opts.Prompter.Select("Select a Platform:", "", platformOptions)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Platform\n", cs.FailureIcon())
				return
			}

			source, err := APICalls.CreateSource(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.SourceName, opts.Platform)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
				return
			}

			sourceId = source.ID
			sourceToken = source.SourceToken
		}

		id := uuid.New()

		cmd := exec.Command("curl",
			"--location",
			"https://in.logfire.ai",
			"--header",
			"Content-Type: application/json",
			"--header",
			fmt.Sprintf("Authorization: Bearer %s", sourceToken),
			"--data",
			fmt.Sprintf("[{\"dt\":\"2023-06-15T6:00:39.351Z\",\"message\":\"%s\"}]", id),
		)
		if err != nil {
			log.Fatal(err)
		}

		go grpcutil.WaitForLog(cfg, id, opts.TeamId, sourceId, stop)

		start := time.Now()

		err = cmd.Run()

		_ = <-stop

		elapsed := time.Since(start)

		fmt.Printf("The round trip took %s\n", elapsed)

	}
}