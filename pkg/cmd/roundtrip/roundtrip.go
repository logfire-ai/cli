package roundtrip

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/helpers"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
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
	SourceId   string
	SourceName string
	Platform   string
	Run        int
	Ctx        context.Context
}

func NewCmdRoundTrip(f *cmdutil.Factory) *cobra.Command {
	opts := &PromptRoundTripOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	opts.Ctx = context.Background()

	cmd := &cobra.Command{
		Use:     "roundtrip",
		Short:   "roundtrip",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			PromptRoundTripRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.TeamId, "team-name", "t", "", "Team name for which the sources will be fetched.")
	cmd.Flags().StringVarP(&opts.SourceId, "source-id", "s", "", "Source ID for which the roundtrip is tested)")
	cmd.Flags().IntVarP(&opts.Run, "run", "r", 0, "Number of rounds")

	return cmd
}

func PromptRoundTripRun(opts *PromptRoundTripOptions) {

	cfg, _ := opts.Config()
	cs := opts.IO.ColorScheme()

	client := http.Client{}

	if opts.TeamId != "" {
		teamId := helpers.TeamNameToTeamId(&client, cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		if teamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s no team with name: %s found.\n", cs.FailureIcon(), opts.TeamId)
			return
		}

		opts.TeamId = teamId
	}

	if opts.TeamId != "" && opts.SourceId != "" {

		source, err := APICalls.GetSource(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.SourceId)
		if err != nil {
			log.Fatal(err)
		}

		id := uuid.New()

		istLocation, err := time.LoadLocation("Asia/Kolkata")
		if err != nil {
			fmt.Println("Error loading IST location:", err)
			return
		}

		// Get the current time in the IST time zone
		currentTime := time.Now().In(istLocation)

		// Format and print the time
		formattedTime := currentTime.Format("2006-01-02 15:04:05")

		ctxCmd, cancelCmd := context.WithTimeout(context.Background(), 5*time.Second)

		cmd := exec.CommandContext(ctxCmd, "curl", "bash", "-c",
			"--location",
			cfg.Get().GrpcIngestion,
			"--header",
			"Content-Type: application/json",
			"--header",
			fmt.Sprintf("Authorization: Bearer %s", source.SourceToken),
			"--header",
			"Diagnostic: True",
			"--header",
			fmt.Sprintf("Github-Run: %v", opts.Run),
			"--data",
			fmt.Sprintf("[{\"dt\":\"%s\",\"message\":\"%s\"}]", formattedTime, id),
		)

		ctx, cancel := context.WithCancel(opts.Ctx)

		go grpcutil.WaitForLog(cfg, id, opts.TeamId, cfg.Get().AccountId, opts.SourceId, ctx, cancel, cancelCmd)

		start := time.Now()

		_ = cmd.Run()

		timeout := 20 * time.Second

		select {
		case <-ctx.Done():
			cancelCmd()
		case <-time.After(time.Until(start.Add(timeout))):
			fmt.Println("Request timed out.")
			os.Exit(1)
		}

		elapsed := time.Since(start)

		fmt.Printf("The round trip took: %s\n", elapsed)
	} else {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		sourceList, err := APICalls.GetAllSources(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err)
			return
		}

		var selectedSource string
		var sourceToken string

		if len(sourceList) != 0 {
			idMap := make(map[string]string)

			var sourceIdNames []string

			for _, source := range sourceList {
				lastFour := ""
				if len(source.ID) > 4 {
					lastFour = source.ID[len(source.ID)-4:]
				} else {
					lastFour = source.ID
				}
				sourceIdNames = append(sourceIdNames, source.Name+" - "+lastFour)
				idMap[source.Name+" - "+lastFour] = source.ID
			}

			selectedSource, err = opts.Prompter.Select("Select Source to round trip:", "", sourceIdNames)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to get sources\n", cs.FailureIcon())
				return
			}

			var ok bool

			opts.SourceId, ok = idMap[selectedSource]
			if !ok {
				log.Fatalf("%s Failed to map to original ID\n\n", cs.FailureIcon())
				return
			}

			for _, source := range sourceList {
				if opts.SourceId == source.ID {
					sourceToken = source.SourceToken
				}
			}
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

			opts.SourceId = source.ID
			sourceToken = source.SourceToken
		}

		id := uuid.New()

		istLocation, err := time.LoadLocation("Asia/Kolkata")
		if err != nil {
			fmt.Println("Error loading IST location:", err)
			return
		}

		// Get the current time in the IST time zone
		currentTime := time.Now().In(istLocation)

		// Format and print the time
		formattedTime := currentTime.Format("2006-01-02 15:04:05")

		ctxCmd, cancelCmd := context.WithTimeout(context.Background(), 10*time.Second)

		cmd := exec.CommandContext(ctxCmd, "curl", "bash", "-c",
			"--location",
			cfg.Get().GrpcIngestion,
			"--header",
			"Content-Type: application/json",
			"--header",
			"Diagnostic: True",
			"--header",
			fmt.Sprintf("Authorization: Bearer %s", sourceToken),
			"--data",
			fmt.Sprintf("[{\"dt\":\"%s\",\"message\":\"%s\"}]", formattedTime, id),
		)

		ctx, cancel := context.WithCancel(opts.Ctx)

		go grpcutil.WaitForLog(cfg, id, opts.TeamId, cfg.Get().AccountId, opts.SourceId, ctx, cancel, cancelCmd)

		start := time.Now()

		_ = cmd.Run()

		timeout := 20 * time.Second

		select {
		case <-ctx.Done():
			cancelCmd()
		case <-time.After(time.Until(start.Add(timeout))):
			fmt.Println("Request timed out.")
			os.Exit(1)
		}

		elapsed := time.Since(start)

		fmt.Printf("The round trip took: %s\n", elapsed)
	}
}
