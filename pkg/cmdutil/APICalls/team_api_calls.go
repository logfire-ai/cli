package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/logfire-sh/cli/pkg/cmd/teams/models"
)

func DeleteTeam(client *http.Client, token string, endpoint string, teamID string) error {
	req, err := http.NewRequest("DELETE", endpoint+"api/team/"+teamID, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Logfire-cli")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var teamDeleteResp models.CreateTeamResponse
	err = json.Unmarshal(body, &teamDeleteResp)
	if err != nil {
		return err
	}

	if !teamDeleteResp.IsSuccessful {
		return errors.New(teamDeleteResp.Message[0])
	}

	return nil
}

func UpdateTeam(client *http.Client, token string, endpoint string, teamID string, teamName string) (models.Team, error) {
	data := models.CreateTeamRequest{
		Name: teamName,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return models.Team{}, err
	}

	req, err := http.NewRequest("PUT", endpoint+"api/team/"+teamID, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.Team{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.Team{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Team{}, err
	}

	var teamUpdateResp models.CreateTeamResponse
	err = json.Unmarshal(body, &teamUpdateResp)
	if err != nil {
		return models.Team{}, err
	}

	if !teamUpdateResp.IsSuccessful {
		return teamUpdateResp.Data, errors.New(teamUpdateResp.Message[0])
	}

	return teamUpdateResp.Data, nil
}

func ListTeams(client *http.Client, token string, endpoint string) ([]models.Team, error) {
	url := endpoint + "api/team"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []models.Team{}, err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return []models.Team{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []models.Team{}, err
	}

	var response models.AllTeamResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.Team{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return []models.Team{}, err
	}

	if !response.IsSuccessful {
		return []models.Team{}, errors.New(response.Message[0])
	}

	return response.Data, nil
}

func CreateTeam(token, endpoint string, teamName string) (models.Team, error) {
	client := http.Client{}

	data := models.CreateTeamRequest{
		Name: teamName,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return models.Team{}, err
	}

	url := endpoint + "api/team"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.Team{}, err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return models.Team{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Team{}, err
	}

	var teamCreateResp models.CreateTeamResponse
	err = json.Unmarshal(body, &teamCreateResp)
	if err != nil {
		return models.Team{}, err
	}

	if !teamCreateResp.IsSuccessful {
		return teamCreateResp.Data, errors.New("failed to create team")
	}

	return teamCreateResp.Data, nil
}
