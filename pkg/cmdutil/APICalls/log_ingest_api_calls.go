package APICalls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"io"
	"net/http"
	"time"
)

type LogMessage struct {
	Dt      string `json:"dt"`
	Message string `json:"message"`
}

func LogIngestFlow(endpoint, sourceToken string) (string, error) {
	istLocation, err := time.LoadLocation("UTC")
	if err != nil {
		return "", fmt.Errorf("Error loading IST location: %v", err)
	}

	currentTime := time.Now().In(istLocation)
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	logMessage := []LogMessage{
		{
			Dt:      formattedTime,
			Message: "Hello from Logfire!",
		},
	}

	reqBody, err := json.Marshal(logMessage)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal request body: %v", err)
	}

	ctxCmd, cancelCmd := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCmd()

	url := endpoint

	transport := &http.Transport{
		IdleConnTimeout:   30 * time.Second,
		MaxIdleConns:      100,
		MaxConnsPerHost:   0,
		DisableKeepAlives: false,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctxCmd, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sourceToken))

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}
		return "", fmt.Errorf("Failed to execute request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
