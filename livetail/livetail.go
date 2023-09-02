package livetail

import (
	"context"
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"
	"net/http"
	"sort"
	"time"

	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/logfire-sh/cli/services/flink-service"
)

type Livetail struct {
	Logs          string
	pbSources     []*pb.Source
	sourcesOffset map[string]uint64
}

var request = &pb.FilterRequest{
	DateTimeFilter:    &pb.DateTimeFilter{},
	FieldBasedFilters: []*pb.FieldBasedFilter{},
	SearchQueries:     []string{},
	Sources:           []*pb.Source{},
	BatchSize:         5,
	IsScrollDown:      true,
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

	}
	return
}

func (l *Livetail) GenerateLogs(stop chan error) {
	filterService := grpcutil.NewFilterService()
	defer filterService.CloseConnection()
	request.Sources = l.pbSources
	for {
		select {
		case <-stop:
			return
		default:
			response, err := filterService.Client.GetFilteredData(context.Background(), request)
			if err != nil {
				stop <- err
				return
			}

			if len(response.Records) > 0 {
				sort.Sort(ByOffset(response.Records))
				l.sourcesOffset = getOffsets(l.sourcesOffset, response.Records)
				l.pbSources = addOffset(l.pbSources, l.sourcesOffset)
				newLogs := showLogsWithColor(response.Records)
				l.Logs += newLogs
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

// Convert logs with colors
func showLogsWithColor(records []*pb.FilteredRecord) string {
	stream := ""
	for _, record := range records {
		stream += fmt.Sprintf("[yellow]" + record.Dt +
			"[green] " + record.SourceName +
			"[blue] " + record.Level + " [white]" +
			record.Message + "\n")
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
func getFilteredData(client pb.FlinkServiceClient, sources []*pb.Source) (*pb.FilteredRecords, error) {
	// Invoke the gRPC method
	request.Sources = sources

	response, err := client.GetFilteredData(context.Background(), request)
	if err != nil {
		return nil, err
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
