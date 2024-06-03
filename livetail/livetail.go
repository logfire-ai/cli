package livetail

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/logfire-sh/cli/services/flink-service"
)

type Livetail struct {
	Logs          string
	pbSources     []*pb.Source
	sourcesOffset map[string]uint64
	offsetMutex   sync.Mutex // Mutex to protect sourcesOffset map
	FilterService *grpcutil.FilterService
}

var request = &pb.FilterRequest{
	DateTimeFilter:    &pb.DateTimeFilter{},
	FieldBasedFilters: []*pb.FieldBasedFilter{},
	SearchQueries:     []string{},
	Sources:           []*pb.Source{},
	BatchSize:         15,
	IsScrollDown:      false,
}

func NewLivetail() (*Livetail, error) {
	livetail := &Livetail{
		Logs:          "",
		sourcesOffset: make(map[string]uint64),
	}

	return livetail, nil
}

// Inverse map
var OperatorToName = map[string]string{
	":":  pb.FieldBasedFilter_Operator_name[0],
	"!:": pb.FieldBasedFilter_Operator_name[1],
	"=":  pb.FieldBasedFilter_Operator_name[2],
	"!=": pb.FieldBasedFilter_Operator_name[3],
	">":  pb.FieldBasedFilter_Operator_name[4],
	">=": pb.FieldBasedFilter_Operator_name[5],
	"<":  pb.FieldBasedFilter_Operator_name[6],
	"<=": pb.FieldBasedFilter_Operator_name[7],
}

func (livetail *Livetail) ApplyFilter(
	cfg config.Config,
	sourceFilter []string,
	StartDateTimeFilter time.Time,
	EndDateTimeFilter time.Time,
	FieldBasedFilterName string,
	FieldBasedFilterValue string,
	FieldBasedFilterCondition string,
) {

	client := &http.Client{}

	var sources []models.Source

	request.AccountID = cfg.Get().AccountId
	request.TeamID = cfg.Get().TeamId

	if sourceFilter != nil {
		for _, sourceId := range sourceFilter {
			source, _ := APICalls.GetSource(cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().TeamId, sourceId)
			sources = append(sources, source)
		}
	} else {
		sources, _ = APICalls.GetAllSources(client, cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().TeamId)
	}

	livetail.pbSources = createGrpcSource(sources)

	if StartDateTimeFilter.IsZero() {
		request.DateTimeFilter.StartTimeStamp = timestamppb.New(time.Now().Add(-1 * time.Second))
	}

	if !StartDateTimeFilter.IsZero() {
		request.DateTimeFilter.StartTimeStamp = timestamppb.New(StartDateTimeFilter)

		if !EndDateTimeFilter.IsZero() {
			request.DateTimeFilter.EndTimeStamp = timestamppb.New(EndDateTimeFilter)
		}
	}

	if FieldBasedFilterName != "" && FieldBasedFilterValue != "" && FieldBasedFilterCondition != "" {
		request.FieldBasedFilters = append(request.FieldBasedFilters, &pb.FieldBasedFilter{
			FieldName:  FieldBasedFilterName,
			FieldValue: FieldBasedFilterValue,
			Operator:   pb.FieldBasedFilter_Operator(pb.FieldBasedFilter_Operator_value[OperatorToName[FieldBasedFilterCondition]]),
		})
	} else {
		request.FieldBasedFilters = []*pb.FieldBasedFilter{}
	}

}

func (l *Livetail) CreateConnection() {
	l.FilterService = grpcutil.NewFilterService()
}

func (l *Livetail) GenerateLogs(ctx context.Context, cfg config.Config) {
	request.Sources = l.pbSources
	theme := cfg.Get().Theme

	for {
		select {
		case <-ctx.Done():
			return
		default:
			response, err := l.FilterService.Client.GetFilteredData(context.Background(), request)

			if err != nil {
				_, cancel := context.WithCancel(ctx)
				defer cancel()

				return
			}

			if len(response.Records) > 0 {
				sort.Sort(ByOffset(response.Records))

				// Lock the mutex before accessing sourcesOffset
				l.offsetMutex.Lock()
				l.sourcesOffset = getOffsets(l.sourcesOffset, response.Records)
				l.offsetMutex.Unlock()

				// Lock the mutex before accessing pbSources
				l.offsetMutex.Lock()
				l.pbSources = addOffset(l.pbSources, l.sourcesOffset)
				l.offsetMutex.Unlock()

				newLogs := showLogsWithColor(response.Records, theme)
				l.Logs += newLogs
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

// Convert logs with colors
func showLogsWithColor(records []*pb.FilteredRecord, theme string) string {
	stream := ""
	if theme == "dark" {
		for _, record := range records {
			stream += fmt.Sprintf("[yellow]" + record.Dt +
				"[green] " + record.SourceName +
				"[blue] " + record.Level + " [white]" +
				record.Message + "\n")
		}

	} else {
		for _, record := range records {
			stream += fmt.Sprintf(`[gray]%s [purple]%s [blue]%s [black]%s`+"\n",
				record.Dt, record.SourceName, record.Level, record.Message)
		}
	}

	return stream
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
// func getFilteredData(client pb.FilterServiceClient, sources []*pb.Source) (*pb.FilteredRecords, error) {
// 	// Invoke the gRPC method
// 	request.Sources = sources

// 	response, err := client.GetFilteredData(context.Background(), request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

func getOffsets(offsets map[string]uint64, records []*pb.FilteredRecord) map[string]uint64 {
	for _, record := range records {
		if offsets[record.SourceName] == 0 || record.Offset >= offsets[record.SourceName] {
			offsets[record.SourceName] = record.Offset + 1
		}
	}
	return offsets
}
