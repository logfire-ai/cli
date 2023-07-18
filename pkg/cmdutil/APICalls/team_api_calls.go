package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmd/teams/models"
	"io"
	"net/http"
)

func DeleteTeam(client *http.Client, token string, teamID string) error {
	req, err := http.NewRequest("DELETE", "https://api.logfire.sh/api/team/"+teamID, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
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
		fmt.Print(teamDeleteResp)
		return errors.New("failed to delete team")
	}

	return nil
}

func UpdateTeam(client *http.Client, token string, teamID string, teamName string) (models.Team, error) {
	data := models.CreateTeamRequest{
		Name: teamName,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return models.Team{}, err
	}

	req, err := http.NewRequest("PUT", "https://api.logfire.sh/api/team/"+teamID, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.Team{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
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
		return teamUpdateResp.Data, errors.New("failed to update team")
	}

	return teamUpdateResp.Data, nil
}
