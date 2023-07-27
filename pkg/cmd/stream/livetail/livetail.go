package livetail

import (
	"context"
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/filters"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LivetailOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	TeamId                    string
	StartDateTimeFilter       string
	EndDateTimeFilter         string
	SourceFilter              []string
	SearchFilter              []string
	FieldBasedFilterName      string
	FieldBasedFilterValue     string
	FieldBasedFilterCondition string
	SaveView                  bool
	ViewName                  string
}

func NewLivetailCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LivetailOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "livetail",
		Short: "Show livetail ",
		Long: heredoc.Docf(`
			Get live stream of logs coming from multiple sources.
		`),
		Example: heredoc.Doc(`
			# start stream of logs
			$ logfire stream livetail --team-id <team-id> --source-id <source-id> --search <search>
			  --field-name <field-name> --field-value <field-value> --field-condition <field-condition> 
			  --start-date <start-date> --end-date <end-date> --save-view <true|default=false> --view-name <view-name>

			# start interactive setup
			$ logfire stream livetail
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			livetailRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team ID for which the sources will be fetched.")
	cmd.Flags().StringSliceVarP(&opts.SourceFilter, "source-id", "s", nil, "Filter logs by sources. (Multiple sources can be specified)")
	cmd.Flags().StringSliceVarP(&opts.SearchFilter, "search", "q", nil, "Filter logs by search.  (Multiple search queries can be specified)")
	cmd.Flags().StringVarP(&opts.FieldBasedFilterName, "field-name", "n", "", "Filter logs by Fields Name (Name, Value, Condition must be specified)")
	cmd.Flags().StringVarP(&opts.FieldBasedFilterValue, "field-value", "v", "", "Filter logs by Fields Value (Name, Value, Condition must be specified)")
	cmd.Flags().StringVarP(&opts.FieldBasedFilterCondition, "field-condition", "c", "", "Filter logs by Fields condition (Name, Value, Condition must be specified)")
	cmd.Flags().StringVarP(&opts.StartDateTimeFilter, "start-date", "", "", "Filter logs by Start date (Example: --start-date now-2d) should give logs from 2Days ago.")
	cmd.Flags().StringVarP(&opts.EndDateTimeFilter, "end-date", "e", "", "Filter logs by End date (Start date should be specified).")
	cmd.Flags().BoolVarP(&opts.SaveView, "save-view", "", false, "Do you want to save the filters as a View. (Default: false)")
	cmd.Flags().StringVarP(&opts.ViewName, "view-name", "", "", "Enter a name for the view.")

	return cmd
}

func livetailRun(opts *LivetailOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Interactive && opts.TeamId == "" && opts.SourceFilter == nil && opts.SearchFilter == nil && opts.FieldBasedFilterName == "" &&
		opts.FieldBasedFilterValue == "" && opts.FieldBasedFilterCondition == "" && opts.StartDateTimeFilter == "" && opts.EndDateTimeFilter == "" {

		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		filterChoice, _ := opts.Prompter.Confirm("Do you want apply any filter?", false)

		if filterChoice {
			filterBySource, _ := opts.Prompter.Confirm("Do you want to filter by source?", false)

			if filterBySource {
				opts.SourceFilter, _ = pre_defined_prompters.AskSourceIds(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)
			}

			filterBySearch, _ := opts.Prompter.Confirm("Do you want to filter by Text search? (You can enter multiple separate words separated by a comma)", false)

			if filterBySearch {
				search, _ := opts.Prompter.Input("Enter the words to filter: (You can enter multiple separate words separated by a comma)", "")

				opts.SearchFilter = append(opts.SearchFilter, search)
			}

			filterByField, _ := opts.Prompter.Confirm("Do you want to filter by Field? (eg: level = info, message [contains] word)", false)

			if filterByField {
				var conditionOptions = []string{
					"CONTAINS",
					"DOES_NOT_CONTAIN",
					"EQUALS",
					"NOT_EQUALS",
					"GREATER_THAN",
					"GREATER_THAN_EQUALS",
					"LESS_THAN",
					"LESS_THAN_EQUALS",
				}

				opts.FieldBasedFilterName, _ = opts.Prompter.Input("Enter the field name:", "")

				opts.FieldBasedFilterCondition, _ = opts.Prompter.Select("Select a condition to match field against value:", "", conditionOptions)

				opts.FieldBasedFilterValue, _ = opts.Prompter.Input("Enter the field value:", "")
			}

			filterByDate, _ := opts.Prompter.Confirm("Do you want to filter by Date?", false)

			if filterByDate {
				opts.StartDateTimeFilter, _ = opts.Prompter.Input("Enter start date: (eg: now-2d = two days ago from now)", "")

				opts.EndDateTimeFilter, _ = opts.Prompter.Input("Enter end date: (Can be left empty) (eg: now-2d = two days ago from now)", "")
			}
		}

		opts.SaveView, _ = opts.Prompter.Confirm("Do you want to save the filters as a view?", false)
		if opts.SaveView {
			opts.ViewName, _ = opts.Prompter.Input("Enter a name for the view:", "")
		}

	} else {
		if opts.TeamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s team-id is required.\n", cs.FailureIcon())
			return
		}
	}

	var sources []models.Source

	if opts.SourceFilter == nil {
		sources, err = APICalls.GetAllSources(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
			return
		}
	} else {
		for _, sourceId := range opts.SourceFilter {
			source, err := APICalls.GetSource(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, sourceId)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
				return
			}
			sources = append(sources, source)
		}
	}

	if opts.SaveView == true {
		err := APICalls.CreateView(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, sources, opts.SearchFilter,
			opts.FieldBasedFilterName, opts.FieldBasedFilterValue, opts.FieldBasedFilterCondition,
			opts.StartDateTimeFilter, opts.EndDateTimeFilter, opts.ViewName)
		if err != nil {
			return
		}
	}

	pbSources := createGrpcSource(sources)
	var sourcesOffset = make(map[string]uint64)

	for {
		response, err := makeGrpcCall(pbSources, opts)
		if err != nil {
			continue
		}

		if len(response.Records) > 0 {
			sort.Sort(ByOffset(response.Records))
			sourcesOffset = getOffsets(sourcesOffset, response.Records)
			pbSources = addOffset(pbSources, sourcesOffset)
			showLogs(opts.IO, response.Records)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// Convert logs with colors
func showLogs(io *iostreams.IOStreams, records []*pb.FilteredRecord) {
	cs := io.ColorScheme()
	for _, record := range records {

		fmt.Fprintf(io.Out, "%s %s [%s] %s\n",
			cs.Yellow(record.Dt), cs.Green(record.SourceName), (cs.Cyan(strings.ToUpper(record.Level))), record.Message)
	}
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

// ByOffset for sorting the records based on Offset
type ByOffset []*pb.FilteredRecord

func (a ByOffset) Len() int           { return len(a) }
func (a ByOffset) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOffset) Less(i, j int) bool { return a[i].Offset < a[j].Offset }

// addOffset adds offset to sources based on the last response received.
func addOffset(sources []*pb.Source, offset map[string]uint64) []*pb.Source {
	for _, source := range sources {
		source.StartingOffset = offset[source.SourceName]
	}

	return sources
}

// getFilteredData makes the actual grpc call to connect with flink-service.
func getFilteredData(client pb.FlinkServiceClient, sources []*pb.Source, opts *LivetailOptions) (*pb.FilteredRecords, error) {
	// Prepare the request payload
	request := &pb.FilterRequest{
		DateTimeFilter:    &pb.DateTimeFilter{},
		FieldBasedFilters: []*pb.FieldBasedFilter{},
		SearchQueries:     []string{},
		Sources:           sources,
		BatchSize:         100,
		IsScrollDown:      true,
	}

	if opts.StartDateTimeFilter == "" {
		request.DateTimeFilter.StartTimeStamp = timestamppb.New(time.Now().Add(-1 * time.Second))
	}

	if opts.StartDateTimeFilter != "" {
		request.DateTimeFilter.StartTimeStamp = timestamppb.New(filters.ShortDateTimeToGoDate(opts.StartDateTimeFilter))

		if opts.EndDateTimeFilter != "" {
			request.DateTimeFilter.EndTimeStamp = timestamppb.New(filters.ShortDateTimeToGoDate(opts.EndDateTimeFilter))
		}
	}

	if opts.SearchFilter != nil {
		request.SearchQueries = append(request.SearchQueries, opts.SearchFilter...)
	}

	if opts.FieldBasedFilterName != "" && opts.FieldBasedFilterValue != "" && opts.FieldBasedFilterCondition != "" {
		request.FieldBasedFilters = append(request.FieldBasedFilters, &pb.FieldBasedFilter{
			FieldName:  opts.FieldBasedFilterName,
			FieldValue: opts.FieldBasedFilterValue,
			Operator:   pb.FieldBasedFilter_Operator(pb.FieldBasedFilter_Operator_value[opts.FieldBasedFilterCondition]),
		})
	}

	// Invoke the gRPC method
	response, err := client.GetFilteredData(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// MakeGrpcCall makes creates a connection and makes a call to the server
func makeGrpcCall(pbSources []*pb.Source, opts *LivetailOptions) (*pb.FilteredRecords, error) {
	grpc_url := "api-stg.logfire.ai:443"
	conn, err := grpc.Dial(grpc_url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewFlinkServiceClient(conn)

	response, err := getFilteredData(client, pbSources, opts)
	if err != nil {
		fmt.Println(err)
		return response, errors.Wrap(err, "[MakeGrpcCall][getFilteredData]")
	}

	return response, nil
}

func getOffsets(offsets map[string]uint64, records []*pb.FilteredRecord) map[string]uint64 {
	for _, record := range records {
		if offsets[record.SourceName] == 0 || record.Offset >= offsets[record.SourceName] {
			offsets[record.SourceName] = record.Offset + 1
		}
	}
	return offsets
}
