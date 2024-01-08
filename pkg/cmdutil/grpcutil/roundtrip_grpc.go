package grpcutil

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func WaitForLog(cfg config.Config, id uuid.UUID, teamId string, accountId string, sourceId string, ctx context.Context, cancel context.CancelFunc, cancelCmd context.CancelFunc) {
	var request = &pb.FilterRequest{
		DateTimeFilter: &pb.DateTimeFilter{},
		Sources:        []*pb.Source{},
		TeamID:         teamId,
		AccountID:      accountId,
		BatchSize:      1,
		IsScrollDown:   true,
	}

	var sources []models.Source

	source, err := APICalls.GetSource(cfg.Get().Token, cfg.Get().EndPoint, teamId, sourceId)
	if err != nil {
		return
	}

	sources = append(sources, source)
	pbSources := CreateGrpcSource(sources)

	request.Sources = pbSources

	request.DateTimeFilter.StartTimeStamp = timestamppb.New(time.Now().Add(-4 * time.Second))

	request.SearchQueries = []string{id.String()}

	filterService := NewFilterService("Diagnostic", "True")
	defer filterService.CloseConnection()

	for {
		select {
		case <-ctx.Done():
			return
		default:

			response, err := filterService.Client.GetFilteredData(context.Background(), request)
			if err != nil {
				log.Printf("Request: %v", request)
				log.Printf("Error %v", err)
				return
			}

			if len(response.Records) > 0 {
				if response.Records[0].Message == id.String() {
					cancelCmd()
					cancel()
				}
			}

			// time.Sleep(500 * time.Millisecond)
		}

	}
}
