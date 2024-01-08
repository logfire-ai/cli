package view

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"

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
			Get stream of logs from specific view.
		`),
		Example: heredoc.Doc(`
			# start stream of logs from specific view
			$ logfire stream view --team-id <team-id> --view-id <view-id>
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
		BatchSize:         15,
		IsScrollDown:      false,
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

	view, err := APICalls.GetView(cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.ViewId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		return
	}

	pbSources := grpcutil.CreateGrpcSource(view.SourcesFilter)
	var sourcesOffset = make(map[string]uint64)

	request.AccountID = cfg.Get().AccountId
	request.TeamID = opts.TeamId
	request.Sources = pbSources
	request.ViewID = view.Id

	filterService := grpcutil.NewFilterService()
	defer filterService.CloseConnection()

	for {
		response, err := filterService.Client.GetFilteredData(context.Background(), request)
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
