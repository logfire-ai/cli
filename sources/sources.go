package sources

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/logfire-sh/cli/models"
)

func GetAllSources(token, teamId, url string) ([]models.Source, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []models.Source{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return []models.Source{}, err
	}
	defer resp.Body.Close()

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

	if sourceResp.IsSuccessful != true {
		fmt.Println("Internal Server Error!")
		return []models.Source{}, errors.New("Api error!")
	}

	return sourceResp.Data, nil
}

func CreateSources(token, teamId, url string, data models.SourceCreate) ([]models.Source, error) {
	client := &http.Client{}

	reqBody, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []models.Source{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return []models.Source{}, err
	}
	defer resp.Body.Close()

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

	if sourceResp.IsSuccessful != true {
		fmt.Println("Internal Server Error!")
		return []models.Source{}, errors.New("Api error!")
	}

	return sourceResp.Data, nil
}

func DeleteSources(token, teamId, url string) ([]models.Source, error) {

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []models.Source{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return []models.Source{}, err
	}
	defer resp.Body.Close()

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

	if sourceResp.IsSuccessful != true {
		fmt.Println("Internal Server Error!")
		return []models.Source{}, errors.New("Api error!")
	}

	return sourceResp.Data, nil
}
