package grpc

import (
	"context"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

func GetLog(cfg config.Config, sourceId string) {
	var sources []models.Source

	source, err := APICalls.GetSource(cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().TeamId, sourceId)
	if err != nil {
		return
	}
	sources = append(sources, source)
	pbSources := createGrpcSource(sources)

	for {
		response, err := makeGrpcCall(pbSources)
		if err != nil {
			continue
		}

		log.Println("hello")

		if len(response.Records) > 0 {
			log.Println(response.Records[0])
		}
		//time.Sleep(500 * time.Millisecond)
	}
}

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

func getFilteredData(client pb.FlinkServiceClient, sources []*pb.Source) (*pb.FilteredRecords, error) {
	// Prepare the request payload
	request := &pb.FilterRequest{
		DateTimeFilter:    &pb.DateTimeFilter{},
		FieldBasedFilters: []*pb.FieldBasedFilter{},
		SearchQueries:     []string{},
		Sources:           sources,
		BatchSize:         100,
		IsScrollDown:      true,
	}

	request.DateTimeFilter.StartTimeStamp = timestamppb.New(time.Now().Add(-1 * time.Second))

	// Invoke the gRPC method
	response, err := client.GetFilteredData(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// MakeGrpcCall makes creates a connection and makes a call to the server
func makeGrpcCall(pbSources []*pb.Source) (*pb.FilteredRecords, error) {
	grpc_url := "api-stg.logfire.ai:443"
	conn, err := grpc.Dial(grpc_url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewFlinkServiceClient(conn)

	response, err := getFilteredData(client, pbSources)
	if err != nil {
		fmt.Println(err)
		return response, errors.Wrap(err, "[MakeGrpcCall][getFilteredData]")
	}

	return response, nil
}
