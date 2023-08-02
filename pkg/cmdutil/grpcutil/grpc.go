package grpcutil

import (
	"context"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

// CreateGrpcSource creates a proper sources to be used in grpc request
func CreateGrpcSource(sources []models.Source) []*pb.Source {
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

// AddOffset adds offset to sources based on the last response received.
func AddOffset(sources []*pb.Source, offset map[string]uint64) []*pb.Source {
	for _, source := range sources {
		source.StartingOffset = offset[source.SourceName]
	}

	return sources
}

func CreateSource(request *pb.Source) {
	cfg, _ := config.NewConfig()
	grpc_url := cfg.Get().GrpcEndpoint
	conn, err := grpc.Dial(grpc_url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewFlinkServiceClient(conn)

	_, err = client.CreateSource(context.Background(), request)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
}

// GetFilteredData makes the actual grpc call to connect with flink-service.
func GetFilteredData(client pb.FlinkServiceClient, request *pb.FilterRequest) (*pb.FilteredRecords, error) {
	// Prepare the request payload

	// Invoke the gRPC method
	response, err := client.GetFilteredData(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// MakeGrpcCall makes creates a connection and makes a call to the server
func MakeGrpcCall(request *pb.FilterRequest) (*pb.FilteredRecords, error) {
	cfg, _ := config.NewConfig()
	grpc_url := cfg.Get().GrpcEndpoint
	conn, err := grpc.Dial(grpc_url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewFlinkServiceClient(conn)

	response, err := GetFilteredData(client, request)
	if err != nil {
		fmt.Println(err)
		return response, errors.Wrap(err, "[MakeGrpcCall][getFilteredData]")
	}

	return response, nil
}

func GetOffsets(offsets map[string]uint64, records []*pb.FilteredRecord) map[string]uint64 {
	for _, record := range records {
		if offsets[record.SourceName] == 0 || record.Offset >= offsets[record.SourceName] {
			offsets[record.SourceName] = record.Offset + 1
		}
	}
	return offsets
}
