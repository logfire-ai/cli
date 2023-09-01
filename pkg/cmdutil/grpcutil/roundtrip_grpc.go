package grpcutil

import (
	"time"

	"github.com/google/uuid"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func WaitForLog(cfg config.Config, id uuid.UUID, teamId string, sourceId string, stop chan bool) {
	var request = &pb.FilterRequest{
		DateTimeFilter: &pb.DateTimeFilter{},
		Sources:        []*pb.Source{},
		BatchSize:      100,
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

	for {
		select {
		case <-stop:
			return
		default:

			response, err := MakeGrpcCall(request)
			if err != nil {
				continue
			}

			if len(response.Records) > 0 {
				for _, r := range response.Records {
					if r.Message == id.String() {
						stop <- true
					}
				}
			}
		}
	}
}
