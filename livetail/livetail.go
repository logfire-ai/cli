package livetail

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	pb "github.com/logfire-sh/cli/services/flink-service"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Livetail struct {
	Logs          string
	pbSources     []*pb.Source
	sourcesOffset map[string]uint64
}

func NewLivetail(token, teamId string, endpoint string) (*Livetail, error) {
	sources, err := getAllSourcesByTeamId(token, teamId, endpoint)
	if err != nil {
		return &Livetail{}, err
	}

	livetail := &Livetail{
		Logs:          "",
		pbSources:     createGrpcSource(sources),
		sourcesOffset: make(map[string]uint64),
	}
	return livetail, nil
}

func (l *Livetail) GenerateLogs() {
	for {
		response, err := makeGrpcCall(l.pbSources)
		if err != nil {
			continue
		}

		if len(response.Records) > 0 {
			sort.Sort(ByOffset(response.Records))
			l.sourcesOffset = getOffsets(l.sourcesOffset, response.Records)
			l.pbSources = addOffset(l.pbSources, l.sourcesOffset)
			new_logs := showLogsWithColor(response.Records)
			l.Logs += new_logs
		}
		time.Sleep(500 * time.Millisecond)
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

func getAllSourcesByTeamId(token, teamId string, endpoint string) ([]models.Source, error) {
	url := endpoint + "api/team/" + teamId + "/source"

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []models.Source{}, err
	}

	// Set the Authorization header with the Bearer token
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return []models.Source{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return []models.Source{}, err
	}

	var sourceResp models.SourcesResponse
	err = json.Unmarshal(body, &sourceResp)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return []models.Source{}, err
	}

	// Check if it is a successful response
	if !sourceResp.IsSuccessful {
		fmt.Println(sourceResp.Message)
		return []models.Source{}, errors.New("Api error!")
	}
	return sourceResp.Data, nil
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

// for sorting the records based on Offset
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
	// Prepare the request payload
	request := &pb.FilterRequest{
		DateTimeFilter: &pb.DateTimeFilter{
			StartTimeStamp: timestamppb.New(time.Now().Add(-5 * time.Second)),
		},
		// FieldBasedFilters: []*pb.FieldBasedFilter{
		// 	{}, // Adjust or populate the field-based filters if needed
		// },
		Sources:      sources,
		BatchSize:    100,
		IsScrollDown: true,
	}

	// Invoke the gRPC method
	response, err := client.GetFilteredData(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// MakeGrpcCall makes creates a connection and makes a call to the server
func makeGrpcCall(pbSources []*pb.Source) (*pb.FilteredRecords, error) {
	grpc_url := fmt.Sprintf("api.logfire.ai:443")
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

func getOffsets(offsets map[string]uint64, records []*pb.FilteredRecord) map[string]uint64 {
	for _, record := range records {
		if offsets[record.SourceName] == 0 || record.Offset >= offsets[record.SourceName] {
			offsets[record.SourceName] = record.Offset + 1
		}
	}
	return offsets
}
