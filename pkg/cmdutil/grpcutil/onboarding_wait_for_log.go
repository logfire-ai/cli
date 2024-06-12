package grpcutil

import (
	"context"
	"log"
	"time"

	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	pb "github.com/logfire-sh/cli/services/flink-service"
)

type Options struct {
	Ctx context.Context
}

func GetLog(config config.Config, token, endpoint, teamId, accountId, sourceId, sourceToken string, stop chan error) {
	request := &pb.FilterRequest{
		DateTimeFilter: &pb.DateTimeFilter{},
		Sources:        []*pb.Source{},
		TeamID:         teamId,
		AccountID:      accountId,
		BatchSize:      1,
	}

	source, err := APICalls.GetSource(token, endpoint, teamId, sourceId)
	if err != nil {
		stop <- err
		return
	}

	sources := []models.Source{source}
	pbSources := CreateGrpcSource(sources)
	request.Sources = pbSources

	filterService := NewFilterService()
	defer filterService.CloseConnection()

	for {
		select {
		case <-stop:
			stop <- nil
			return
		default:
			_, err := APICalls.LogIngestFlow(config.Get().GrpcIngestion, sourceToken)
			if err != nil {
				log.Printf("Error: %s", err.Error())
				continue
			}

			// wait some time to process log messages
			time.Sleep(1 * time.Second)

			// Check response from filter service
			response, err := filterService.Client.GetFilteredData(context.Background(), request)
			if err != nil {
				log.Println("Error getting filtered data:", err)
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
