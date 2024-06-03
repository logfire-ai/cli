package grpcutil

import (
	"context"
	"log"
	"time"

	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	pb "github.com/logfire-sh/cli/services/flink-service"
)

func GetLog(token string, endpoint string, teamId string, accountId string, sourceId string, stop chan error) {
	request := &pb.FilterRequest{
		DateTimeFilter: &pb.DateTimeFilter{},
		Sources:        []*pb.Source{},
		TeamID:         teamId,
		AccountID:      accountId,
		BatchSize:      1,
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

			log.Printf("Response: %s\n", response.Records)

			if len(response.Records) > 0 {
				stop <- nil
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}
