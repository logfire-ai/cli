package livetail

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"logfire/models"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	pb "logfire/logfire/flink-service"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getAllSourcesByTeamId(token, teamId string) ([]models.Source, error) {
	url := "https://api.logfire.sh/api/team/" + teamId + "/source"

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

	var sourceResp models.SourceResponse
	err = json.Unmarshal(body, &sourceResp)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return []models.Source{}, err
	}

	// Check if it is a successful response
	if !sourceResp.IsSuccessful {
		fmt.Println("Internal Server Error!")
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
		FieldBasedFilters: []*pb.FieldBasedFilter{
			{}, // Adjust or populate the field-based filters if needed
		},
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
	conn, err := grpc.Dial("api.logfire.sh:443", grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
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

func showLogsWithColor(records []*pb.FilteredRecord) {
	for _, record := range records {
		color.New(color.FgCyan).Print(record.Dt)
		color.New(color.FgGreen).Print(" [" + record.SourceName + "]")
		color.New(color.FgBlue).Print(" [" + record.Level + "] ")
		color.New(color.FgWhite).Print(record.Message + "\n")
	}
}

func getOffsets(offsets map[string]uint64, records []*pb.FilteredRecord) map[string]uint64 {
	for _, record := range records {
		if offsets[record.SourceName] == 0 || record.Offset >= offsets[record.SourceName] {
			offsets[record.SourceName] = record.Offset + 1
		}
	}
	return offsets
}

func gracefulShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine to handle the interrupt signal
	go func() {
		<-signalChan
		color.Red("Closing livetail.")
		os.Exit(0)
	}()
}

func ShowLivetail(token, teamID string) error {
	color.Green("Starting livetail.")

	gracefulShutdown()

	sources, err := getAllSourcesByTeamId(token, teamID)
	if err != nil {
		return errors.Wrap(err, "Error while getting sources!!!")
	}

	pbSources := createGrpcSource(sources)
	sourcesOffset := make(map[string]uint64)

	for {
		response, err := makeGrpcCall(pbSources)
		if err != nil {
			log.Println("Err while getting logs!")
			continue
		}

		if len(response.Records) > 0 {
			sort.Sort(ByOffset(response.Records))
			sourcesOffset = getOffsets(sourcesOffset, response.Records)
			pbSources = addOffset(pbSources, sourcesOffset)
			showLogsWithColor(response.Records)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
