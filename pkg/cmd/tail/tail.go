package tail

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/logfire-sh/cli/pkg/cmdutil/helpers"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"

	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"

	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/filters"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TailOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	Interacted  bool

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
	GUI                       bool
}

func NewTailCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &TailOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "tail",
		Short: "Show tail ",
		Long: heredoc.Docf(`
			Get live stream of logs coming from multiple sources.
		`),
		Example: heredoc.Doc(`
			# start stream of logs
			$ logfire stream tail --team-name <team-name> --source-id <source-id> --search <search>
			  --field-name <field-name> --field-value <field-value> --field-condition <field-condition>
			  --start-date <start-date> --end-date <end-date> --save-view <true|default=false> --view-name <view-name>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			tailRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.TeamId, "team-name", "t", "", "Team name for which the sources will be fetched.")
	cmd.Flags().StringSliceVarP(&opts.SourceFilter, "source-id", "s", nil, "Filter logs by sources. (Multiple sources can be specified)")
	cmd.Flags().StringSliceVarP(&opts.SearchFilter, "search", "q", nil, "Filter logs by search.  (Multiple search queries can be specified)")
	cmd.Flags().StringVarP(&opts.FieldBasedFilterName, "field-name", "n", "", "Filter logs by Fields Name (Name, Value, Condition must be specified)")
	cmd.Flags().StringVarP(&opts.FieldBasedFilterValue, "field-value", "v", "", "Filter logs by Fields Value (Name, Value, Condition must be specified)")
	cmd.Flags().StringVarP(&opts.FieldBasedFilterCondition, "field-condition", "c", "", "Filter logs by Fields condition (Name, Value, Condition must be specified)")
	cmd.Flags().StringVarP(&opts.StartDateTimeFilter, "start-date", "", "", "Filter logs by Start date (Example: --start-date now-2d) should give logs from 2Days ago.")
	cmd.Flags().StringVarP(&opts.EndDateTimeFilter, "end-date", "e", "", "Filter logs by End date (Start date should be specified).")
	cmd.Flags().BoolVarP(&opts.SaveView, "save-view", "", false, "Do you want to save the filters as a View. (Default: false)")
	cmd.Flags().StringVarP(&opts.ViewName, "view-name", "", "", "Enter a name for the view.")
	cmd.Flags().BoolVarP(&opts.GUI, "gui", "", false, "Enable GUI.")

	return cmd
}

func tailRun(opts *TailOptions) {
	conditions := []string{"CONTAINS", "DOES_NOT_CONTAIN", "EQUALS", "NOT_EQUALS", "GREATER_THAN", "GREATER_THAN_EQUALS", "LESS_THAN", "LESS_THAN_EQUALS"}

	var request = &pb.FilterRequest{
		DateTimeFilter:    &pb.DateTimeFilter{},
		FieldBasedFilters: []*pb.FieldBasedFilter{},
		SearchQueries:     []string{},
		Sources:           []*pb.Source{},
		BatchSize:         15,
		IsScrollDown:      false,
	}

	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	client := http.Client{}

	if opts.Interactive && opts.TeamId == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		opts.Interacted = true
	} else {
		if opts.TeamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s team-name is required.\n", cs.FailureIcon())
			return
		}
	}

	if opts.TeamId != "" && !opts.Interacted {
		teamId := helpers.TeamNameToTeamId(&client, cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		if teamId == "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s no team with name: %s found.\n", cs.FailureIcon(), opts.TeamId)
			return
		}

		opts.TeamId = teamId
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
			source, err := APICalls.GetSource(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, sourceId)
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
				return
			}
			sources = append(sources, source)
		}
	}

	if opts.SaveView {
		err := APICalls.CreateView(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, sources, opts.SearchFilter,
			opts.FieldBasedFilterName, opts.FieldBasedFilterValue, opts.FieldBasedFilterCondition,
			opts.StartDateTimeFilter, opts.EndDateTimeFilter, opts.ViewName)
		if err != nil {
			return
		}
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

	if opts.FieldBasedFilterName != "" && opts.FieldBasedFilterValue != "" && opts.
		FieldBasedFilterCondition != "" && helpers.StringNotInArray(strings.ToUpper(opts.FieldBasedFilterCondition), conditions) {
		request.FieldBasedFilters = append(request.FieldBasedFilters, &pb.FieldBasedFilter{
			FieldName:  opts.FieldBasedFilterName,
			FieldValue: opts.FieldBasedFilterValue,
			Operator:   pb.FieldBasedFilter_Operator(pb.FieldBasedFilter_Operator_value[strings.ToUpper(opts.FieldBasedFilterCondition)]),
		})
	}

	pbSources := grpcutil.CreateGrpcSource(sources)
	var sourcesOffset = make(map[string]uint64)

	request.Sources = pbSources
	request.AccountID = cfg.Get().AccountId
	request.TeamID = opts.TeamId

	filterService := grpcutil.NewFilterService()
	defer filterService.CloseConnection()

	for {
		response, err := filterService.Client.GetFilteredData(context.Background(), request)
		if err != nil {
			//log.Fatal(err)
			continue
		}

		if len(response.Records) > 0 {
			sort.Sort(grpcutil.ByOffset(response.Records))
			sourcesOffset = grpcutil.GetOffsets(sourcesOffset, response.Records)
			pbSources = grpcutil.AddOffset(pbSources, sourcesOffset)
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
			cs.Yellow(record.Dt), cs.Green(record.SourceName), cs.Cyan(strings.ToUpper(record.Level)), record.Message)
	}
}
