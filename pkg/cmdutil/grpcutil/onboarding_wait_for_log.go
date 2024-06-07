package grpcutil

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
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
			istLocation, err := time.LoadLocation("UTC")
			if err != nil {
				log.Println("Error loading IST location:", err)
				stop <- err
				return
			}

			currentTime := time.Now().In(istLocation)
			formattedTime := currentTime.Format("2006-01-02 15:04:05")

			ctxCmd, cancelCmd := context.WithTimeout(context.Background(), 5*time.Second)

			cmd := exec.CommandContext(ctxCmd, "curl",
				"--location",
				config.Get().GrpcIngestion,
				"--header", "Content-Type: application/json",
				"--header", fmt.Sprintf("Authorization: Bearer %s", sourceToken),
				"--data", fmt.Sprintf("[{\"dt\":\"%s\",\"message\":\"%s\"}]", formattedTime, "Hello from Logfire!"),
			)

			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			// Start the curl command
			if err := cmd.Start(); err != nil {
				log.Println("Error starting curl command:", err)
				cancelCmd()
				stop <- err
				return
			}

			// Allow the command to run for a short period before cancelling it
			time.Sleep(1100 * time.Millisecond)
			cancelCmd()

			// Wait for the curl command to complete
			if err := cmd.Wait(); err != nil && err.Error() != "signal: killed" {
				log.Println("Error waiting for curl command:", err)
				log.Println("Curl output:", out.String())
				stop <- err
				return
			}

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
