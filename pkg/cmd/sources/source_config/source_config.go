package source_config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type SourceConfigurationOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	TeamId   string
	SourceId string
}

func NewSourceConfigCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &SourceConfigurationOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "configuration",
		Args:  cobra.ExactArgs(0),
		Short: "Get source configuration",
		Long: heredoc.Docf(`
		Get source source for a particular source.
		`),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sources configuration

			# start argument setup
			$ logfire sources configuration --team-id <team-id>  --source-id <source-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			sourceListRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.SourceId, "source-id", "", "Source ID for which the source is to be deleted.")
	cmd.Flags().StringVar(&opts.TeamId, "team-id", "", "Team ID for which the sources will be fetched.")
	return cmd
}

func sourceListRun(opts *SourceConfigurationOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Interactive && opts.TeamId == "" && opts.SourceId == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)
	} else {
		if opts.TeamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s team-id is required.\n", cs.FailureIcon())
			return
		}

		if opts.SourceId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s source-id is required.\n", cs.FailureIcon())
			return
		}
	}

	opts.SourceId, _ = pre_defined_prompters.AskSourceId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

	source, err := APICalls.GetSource(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.SourceId)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Fprintf(opts.IO.ErrOut, "\n%s Error: Connection failed (Server down or no internet)\n", cs.FailureIcon())
			return
		}
		fmt.Fprintf(opts.IO.ErrOut, "\n%s Failed to load sources\n", cs.FailureIcon())
		return
	}

	configuration, err := APICalls.GetConfiguration(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.SourceId)
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(configuration)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println(err)
		return
	}

	printKeys("", data, source.ID, source.SourceToken, opts)
}

func promptForOS(opts *SourceConfigurationOptions) string {

	os, err := opts.Prompter.Select("Pick Your OS for Configuration", "", []string{"CentOS", "Ubuntu", "Windows", "MacOS", "Other"})
	if err != nil {
		log.Fatal(err)
	}

	return strings.ToLower(os)
}

func printKeys(prefix string, data interface{}, id, token string, opts *SourceConfigurationOptions) {
	//HR
	rulerColor := color.New(color.FgHiBlack).Add(color.Bold)

	// Heading
	headingColor := color.New(color.FgWhite, color.Bold)

	// Step
	stepColor := color.New(color.FgWhite, color.Bold)

	// Title
	titleColor := color.New(color.FgBlue)

	// Code
	codeColor := color.New(color.FgWhite)

	switch v := data.(type) {
	case map[string]interface{}:
		var keys []string
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		// Special case: if the map contains both "title" and "code", print them in that order
		if _, titleExists := v["title"]; titleExists {
			if _, codeExists := v["code"]; codeExists {

				title := v["title"].(string)

				if strings.Contains(title, "#{") {
					replacedLink := strings.ReplaceAll(title, "#{", "")
					replacedLink = strings.ReplaceAll(replacedLink, "}", "")

					titleColor.Printf("      %s\n", replacedLink)
				} else {
					titleColor.Printf("      %s\n", v["title"])
				}

				unformattedCode := v["code"]

				code := strings.Split(unformattedCode.(string), "\n")

				for _, c := range code {
					if strings.HasPrefix(c, "//") {
						rulerColor.Printf("         %s\n", c)
					} else {
						var parsedCode string
						parsedCode = strings.ReplaceAll(c, "&{source_id}", id)
						parsedCode = strings.ReplaceAll(parsedCode, "${source_token}", token)

						codeColor.Printf("         %s\n", parsedCode)
					}
				}

				return
			}
		}

		for _, key := range keys {
			value := v[key]
			keyIsNum, err := strconv.Atoi(key)
			if err == nil {
				// The key is a number (Step)
				if keyIsNum == 1 {
					// Special handling for Step 1
					selectedOS := promptForOS(opts)
					if osSteps, ok := value.(map[string]interface{}); ok {
						if osStep, exists := osSteps[selectedOS]; exists {
							stepColor.Printf("\n    %sStep %d (%s):\n", prefix, keyIsNum, selectedOS)
							printKeys(prefix+"  ", osStep, id, token, opts)
							continue
						}
					}
				} else {
					stepColor.Printf("\n    %sStep %d:\n", prefix, keyIsNum)
				}
			} else {
				// The key is a string
				if key != "code" && key != "title" {
					rulerColor.Print("****************************************************************************************")
					headingColor.Printf("\n%s:\n", key)
				}
			}

			// Recur with a new prefix for steps; otherwise, keep it the same
			if err == nil {
				printKeys(prefix+"  ", value, id, token, opts)
			} else {
				printKeys(prefix, value, id, token, opts)
			}
		}
	case []interface{}:
		for _, value := range v {
			printKeys(prefix, value, id, token, opts)
		}
	default:
		fmt.Printf("%s%v\n", prefix, v)
	}
}
