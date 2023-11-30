package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	sourceModels "github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmd/sql/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type SQLQueryOptions struct {
	IO *iostreams.IOStreams

	HttpClient func() *http.Client
	Prompter   prompter.Prompter
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	SQLQuery    string
	Role        string
}

func NewCmdSql(f *cmdutil.Factory) *cobra.Command {
	opts := &SQLQueryOptions{
		IO: f.IOStreams,

		HttpClient: f.HttpClient,
		Prompter:   f.Prompter,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "sql",
		Short: "Run a sql query",
		Long:  "Run a sql query",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sql

			# start argument setup
			$ logfire sql --team-id <team-id> --query <query>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			SqlQueryRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team id to be queried.")
	cmd.Flags().StringVarP(&opts.SQLQuery, "sql-query", "q", "", "SQL Query.")

	return cmd
}

func GetRecommendations(opts *SQLQueryOptions, cfg config.Config) {
	opts.IO.StartProgressIndicatorWithLabel("Getting recommendations, please wait...")

	recommendations, err := APICalls.GetRecommendations(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, cfg.Get().Role)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s",err)
	}

	var options []string

	for _, recommendation := range recommendations.Data {
		options = append(options, fmt.Sprintf(`Title: %s`, recommendation.CaptionTitle))
	}

	opts.IO.StopProgressIndicator()

	selectedQuery, _ := opts.Prompter.Select("Select a recommended query to run", "", options)

	var selectedQueryParsed string

	ret := regexp.MustCompile(`Title:\s+(.+)`)

	for _, recommendation := range recommendations.Data {
		submatches := ret.FindStringSubmatch(selectedQuery)
		if recommendation.CaptionTitle == submatches[1] {
			selectedQueryParsed = recommendation.SQLStatement
		}
	}

	if len(selectedQueryParsed) > 1 {
		opts.SQLQuery = selectedQueryParsed
		fmt.Fprintf(opts.IO.Out, "Query: %s\n", selectedQueryParsed)
	} else {
		log.Fatal("Query not found in the input.")
	}
}

func SqlQueryRun(opts *SQLQueryOptions) {

	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive && opts.TeamId == "" && opts.SQLQuery == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		choices := []string{
			"Receive AI-generated query recommendations.",
			"Manually enter your query.",
		}

		chooseMode, _ := opts.Prompter.Select("", "", choices)

		if chooseMode == "Receive AI-generated query recommendations." {
			GetRecommendations(opts, cfg)
		} else if chooseMode == "Manually enter your query." {
			opts.SQLQuery, _ = opts.Prompter.Input("Write your SQL query:", "")
		}

	} else {
		if opts.TeamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
		}

		if opts.SQLQuery == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s SQL Query is required.\n", cs.FailureIcon())
		}
	}

	var sources []sourceModels.Source

	sources, err = APICalls.GetAllSources(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		return
	}

	pbSources := createGrpcSource(sources)

	// Prepare the request payload
	request := &pb.SQLRequest{
		Sql:       opts.SQLQuery,
		Sources:   pbSources,
		BatchSize: 100,
	}

	filterService := grpcutil.NewFilterService()
	defer filterService.CloseConnection()

	response, err := filterService.Client.SubmitSQL(context.Background(), request)
	if err != nil {
		return
	}

	if len(response.Data) > 0 {
		showQuery(opts.IO, response.Data)
	}
}

// Convert logs with colors
func showQuery(io *iostreams.IOStreams, records string) {
	var parsedData models.SQLResponse
	err := json.Unmarshal([]byte(records), &parsedData)
	if err != nil {
		return
	}

	var fieldsNames []string

	for _, field := range parsedData.Fields {
		fieldsNames = append(fieldsNames, field.Name)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(fieldsNames)
	table.SetAutoWrapText(false)

	for _, record := range parsedData.Records {
		var row []string

		for _, field := range parsedData.Fields {
			name := field.Name
			strValue := fmt.Sprintf("%v", record[name])

			// Truncate if length is more than 150 characters
			if len(strValue) > 150 {
				strValue = strValue[:150] + "..." // Truncate and add ellipsis
			}

			re := regexp.MustCompile(`\s+`)

			row = append(row, re.ReplaceAllString(strValue, " "))
		}

		table.Append(row)

		row = []string{}
	}

	table.SetRowLine(true)

	table.Render()
}

// createGrpcSource creates a proper sources to be used in grpc request
func createGrpcSource(sources []sourceModels.Source) []*pb.Source {
	var grpcSources []*pb.Source
	for _, source := range sources {
		pbSource := pb.Source{
			SourceID:   "source_topic_" + source.ID,
			SourceName: source.Name,
			TeamID:     source.TeamID,
		}
		grpcSources = append(grpcSources, &pbSource)
	}
	return grpcSources
}

// getSQL makes the actual grpc call to connect with flink-service.
func getSQL(client pb.FilterServiceClient, sources []*pb.Source, opts *SQLQueryOptions) (*pb.SQLResponse, error) {
	// Prepare the request payload
	request := &pb.SQLRequest{
		Sql:       opts.SQLQuery,
		Sources:   sources,
		BatchSize: 100,
	}

	// Invoke the gRPC method
	response, err := client.SubmitSQL(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
