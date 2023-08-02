package view

import (
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"github.com/spf13/cobra"
)

type ViewStreamOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	TeamId                    string
	ViewId                    string
	StartDateTimeFilter       *time.Time
	EndDateTimeFilter         *time.Time
	SourceFilter              []string
	SearchFilter              []string
	FieldBasedFilterName      string
	FieldBasedFilterValue     string
	FieldBasedFilterCondition string
}

func NewViewStreamOptionsCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ViewStreamOptions{
		IO:          f.IOStreams,
		Prompter:    f.Prompter,
		HttpClient:  f.HttpClient,
		Config:      f.Config,
		Interactive: false,
	}

	cmd := &cobra.Command{
		Use:   "view",
		Args:  cobra.ExactArgs(0),
		Short: "Stream view",
		Long: heredoc.Docf(`
			Get stream of logs from selected view.
		`),
		Example: heredoc.Doc(`
			# start stream of logs from selected view
			$ logfire stream view --team-id <team-id> --view-id <view-id>

			# start interactive setup
			$ logfire stream view
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			ViewStreamRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team ID for which the sources will be fetched.")
	cmd.Flags().StringVarP(&opts.ViewId, "view-id", "v", "", "Team ID for which the sources will be fetched.")
	return cmd
}

func ViewStreamRun(opts *ViewStreamOptions) {
	var request = &pb.FilterRequest{
		DateTimeFilter:    &pb.DateTimeFilter{},
		FieldBasedFilters: []*pb.FieldBasedFilter{},
		SearchQueries:     []string{},
		Sources:           []*pb.Source{},
		BatchSize:         100,
		IsScrollDown:      true,
	}

	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.TeamId == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s team-id is required.\n", cs.FailureIcon())
		os.Exit(0)
	}

	if opts.ViewId == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s view-id is required.\n", cs.FailureIcon())
		os.Exit(0)
	}

	view, err := APICalls.GetView(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.ViewId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		return
	}

	if !view.DateFilter.StartDate.IsZero() {
		request.DateTimeFilter.StartTimeStamp = timestamppb.New(view.DateFilter.StartDate)

		if !view.DateFilter.EndDate.IsZero() {
			request.DateTimeFilter.EndTimeStamp = timestamppb.New(view.DateFilter.EndDate)
		}
	}

	if len(view.SearchFilter) != 0 {
		for _, v := range view.SearchFilter {
			if v.Key != "" {
				if v.Condition != "" {
					if v.Value != "" {
						request.FieldBasedFilters = append(request.FieldBasedFilters, &pb.FieldBasedFilter{
							FieldName:  v.Key,
							FieldValue: v.Value,
							Operator:   pb.FieldBasedFilter_Operator(pb.FieldBasedFilter_Operator_value[v.Condition]),
						})
					}
				}
			}
		}
	}

	if len(view.TextFilter) != 0 {
		for _, t := range view.TextFilter {
			request.SearchQueries = append(request.SearchQueries, t)
		}
	}

	pbSources := grpcutil.CreateGrpcSource(view.SourcesFilter)
	var sourcesOffset = make(map[string]uint64)

	for {
		response, err := grpcutil.MakeGrpcCall(request)
		if err != nil {
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
