package grpcutil

import (
	"context"

	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	pb "github.com/logfire-sh/cli/services/flink-service"
	"google.golang.org/grpc"

	"log"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type FilterService struct {
	conn   *grpc.ClientConn
	Client pb.FilterServiceClient
}

func authUnaryInterceptor(kv ...string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := metadata.Pairs(kv...)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (fs *FilterService) CloseConnection() {
	err := fs.conn.Close()
	if err != nil {
		return
	}
}

func NewFilterService(kv ...string) *FilterService {
	cfg, _ := config.NewConfig()
	grpc_url := cfg.Get().GrpcEndpoint
	allParams := make([]string, 0, len(kv)+2)
	allParams = append(allParams, "Authorization", "Bearer "+cfg.Get().Token)
	allParams = append(allParams, kv...)

	conn, err := grpc.Dial(grpc_url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), grpc.WithUnaryInterceptor(authUnaryInterceptor(allParams...)), grpc.WithUserAgent("Logfire-cli"))
	// conn, err := grpc.Dial(grpc_url, grpc.WithInsecure(), grpc.WithUnaryInterceptor(authUnaryInterceptor(allParams...)), grpc.WithUserAgent("Logfire-cli"))

	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	// Create a gRPC client
	client := pb.NewFilterServiceClient(conn)

	return &FilterService{
		conn:   conn,
		Client: client,
	}
}

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

func GetOffsets(offsets map[string]uint64, records []*pb.FilteredRecord) map[string]uint64 {
	for _, record := range records {
		if offsets[record.SourceName] == 0 || record.Offset >= offsets[record.SourceName] {
			offsets[record.SourceName] = record.Offset + 1
		}
	}
	return offsets
}
