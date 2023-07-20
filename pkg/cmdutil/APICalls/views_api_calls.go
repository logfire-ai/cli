package APICalls

import (
	"encoding/json"
	"errors"
	"github.com/logfire-sh/cli/pkg/cmd/views/models"
	"io"
	"net/http"
)

func DeleteView(client *http.Client, token string, endpoint string, teamId string, viewId string) error {
	req, err := http.NewRequest("DELETE", endpoint+"api/team/"+teamId+"/view/"+viewId, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

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

	var InviteMemberResp models.DeleteViewResponse
	err = json.Unmarshal(body, &InviteMemberResp)
	if err != nil {
		return err
	}

	if !InviteMemberResp.IsSuccessful {
		return errors.New("failed to delete view")
	}

	return nil
}

func ListView(client *http.Client, token string, endpoint string, teamId string) ([]models.ViewResponse, error) {
	req, err := http.NewRequest("GET", endpoint+"api/team/"+teamId+"/view", nil)
	if err != nil {
		return []models.ViewResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []models.ViewResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.ViewResponse{}, err
	}

	var ListViewResp models.ListViewResponse
	err = json.Unmarshal(body, &ListViewResp)
	if err != nil {
		return []models.ViewResponse{}, err
	}

	if !ListViewResp.IsSuccessful {
		return []models.ViewResponse{}, errors.New("failed to delete view")
	}

	return ListViewResp.Views, err
}
