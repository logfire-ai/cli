package grpcutil

import (
	"context"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func GetLog(token string, endpoint string, teamId string, sourceId string, stop chan error) {
	request := &pb.FilterRequest{
		DateTimeFilter: &pb.DateTimeFilter{},
		Sources:        []*pb.Source{},
		BatchSize:      1,
		IsScrollDown:   true,
	}

	var sources []models.Source

	source, err := APICalls.GetSource(token, endpoint, teamId, sourceId)
	if err != nil {
		stop <- err
		return
	}
	sources = append(sources, source)
	pbSources := CreateGrpcSource(sources)

	request.Sources = pbSources

	request.DateTimeFilter.StartTimeStamp = timestamppb.New(time.Now().Add(-1 * time.Second))

	filterService := NewFilterService()
	defer filterService.CloseConnection()

	for {
		select {
		case <-stop:
			stop <- err
			return
		default:
			response, err := filterService.Client.GetFilteredData(context.Background(), request)
			if err != nil {
				stop <- err
				return
			}

			if len(response.Records) > 0 {
				stop <- nil
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}
