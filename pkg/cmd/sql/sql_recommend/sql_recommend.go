package sql

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/sql/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/spf13/cobra"
)

type Temp struct {
	Prompter prompter.Prompter
}

func NewSQLRecommendCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &models.SQLQueryOptions{
		IO: f.IOStreams,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	temp := &Temp{
		Prompter: f.Prompter,
	}

	cmd := &cobra.Command{
		Use:   "recommend",
		Short: "Get a list of query recommendations",
		Long:  "Get a list of query recommendations",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sql recommend

			# start argument setup
			$ logfire sql recommend --team-id <team-id> --role <role>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			SqlRecommendRun(opts, temp, f)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team id to be queried.")
	cmd.Flags().StringVarP(&opts.Role, "role", "r", "", "Your Role.")
	return cmd
}

func SqlRecommendRun(opts *models.SQLQueryOptions, temp *Temp, f *cmdutil.Factory) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive && opts.TeamId == "" && opts.Role == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, temp.Prompter)

		opts.Role, _ = temp.Prompter.Input("Please enter your role: example (Software Developer)", "")

		if opts.Role == "" {
			fmt.Println("Role is required")
			os.Exit(0)
		}

		recommendations, _ := APICalls.GetRecommendations(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.Role)

		var options []string

		for _, recommendation := range recommendations.Data {
			options = append(options, fmt.Sprintf(`
Title: %s
Description: %s
Query: %s
`, recommendation.CaptionTitle,
				recommendation.CaptionDescription,
				recommendation.SQLStatement,
			))
		}

		selectedQuery, _ := temp.Prompter.Select("Select a recommended query to run", "", options)

		// Define a regular expression to match the Query section
		re := regexp.MustCompile(`Query:\s+(.+)`)

		// Find the submatch (the query) in the input string
		submatches := re.FindStringSubmatch(selectedQuery)

		if len(submatches) > 1 {
			opts.SQLQuery = submatches[1]
		} else {
			log.Fatal("Query not found in the input.")
		}

		cmdutil.SqlQueryRun(opts, f)

	} else {
		if opts.TeamId == "" || opts.Role == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Team id and Role is required.\n", cs.FailureIcon())
		} else {

			recommendations, _ := APICalls.GetRecommendations(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.Role)

			var options []string

			for _, recommendation := range recommendations.Data {
				options = append(options, fmt.Sprintf(`
Title: %s
Description: %s
Query: %s
`, recommendation.CaptionTitle,
					recommendation.CaptionDescription,
					recommendation.SQLStatement,
				))
			}

			fmt.Println(options)

			os.Exit(0)
		}
	}

}
