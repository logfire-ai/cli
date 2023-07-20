package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	SQLModels "github.com/logfire-sh/cli/pkg/cmd/sql/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
	"os"
)

type SqlQueryOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	SQLQuery    string
}

func NewCmdSql(f *cmdutil.Factory) *cobra.Command {
	opts := &SqlQueryOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,

		HttpClient: f.HttpClient,
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

func SqlQueryRun(opts *SqlQueryOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.TeamId == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team id is required.\n", cs.FailureIcon())
	}

	if opts.SQLQuery == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s SQL Query is required.\n", cs.FailureIcon())
	}

	var sources []models.Source

	sources, err = APICalls.GetAllSources(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		return
	}

	pbSources := createGrpcSource(sources)

	response, err := makeGrpcCall(pbSources, opts)
	if err != nil {
		return
	}

	if len(response.Data) > 0 {
		showQuery(opts.IO, response.Data)
	}
}

// Convert logs with colors
func showQuery(io *iostreams.IOStreams, records string) {
	var parsedData SQLModels.SQLResponse
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

	for _, record := range parsedData.Records {
		var row []string

		for _, field := range parsedData.Fields {
			name := field.Name
			row = append(row, record[name].(string))
		}

		table.Append(row)

		row = []string{}
	}

	table.SetRowLine(true)

	table.Render()
}

// createGrpcSource creates a proper sources to be used in grpc request
func createGrpcSource(sources []models.Source) []*pb.Source {
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
func getSQL(client pb.FlinkServiceClient, sources []*pb.Source, opts *SqlQueryOptions) (*pb.SQLResponse, error) {
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

// MakeGrpcCall makes creates a connection and makes a call to the server
func makeGrpcCall(pbSources []*pb.Source, opts *SqlQueryOptions) (*pb.SQLResponse, error) {
	grpc_url := "api.logfire.sh:443"
	conn, err := grpc.Dial(grpc_url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewFlinkServiceClient(conn)

	response, err := getSQL(client, pbSources, opts)
	if err != nil {
		fmt.Println(err)
		return response, errors.Wrap(err, "[MakeGrpcCall][getSQL]")
	}

	return response, nil
}