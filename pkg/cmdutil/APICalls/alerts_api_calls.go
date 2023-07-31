package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/models"
	"io"
	"net/http"
)

func CreateAlert(client *http.Client, token string, endpoint string, teamId string, name string, viewId string, numberOfRecords uint32,
	withInSeconds uint32, IntegrationsId []string) error {

	integrationsList, err := GetAlertIntegrations(client, token, endpoint, teamId)
	if err != nil {
		return nil
	}

	var integrationParsed []models.AlertIntegrationBody

	for _, integration := range integrationsList {
		for _, integrationID := range IntegrationsId {
			if integrationID == integration.ModelId {
				integrationParsed = append(integrationParsed, integration)
			}
		}
	}

	data := models.CreateAlertRequest{
		Name:            name,
		ViewId:          viewId,
		NumberOfRecords: numberOfRecords,
		WithinSeconds:   withInSeconds,
		Integrations:    integrationParsed,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint+"api/team/"+teamId+"/alert", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "Logfire-cli")
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

	var CreateAlertResp models.CreateAlertResponse
	err = json.Unmarshal(body, &CreateAlertResp)
	if err != nil {
		return err
	}

	if !CreateAlertResp.IsSuccessful {
		return errors.New("failed to create alert")
	}

	return nil
}

func ListAlert(client *http.Client, token string, endpoint string, teamId string) ([]models.CreateAlertBody, error) {
	req, err := http.NewRequest("GET", endpoint+"api/team/"+teamId+"/alert", nil)
	if err != nil {
		return []models.CreateAlertBody{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []models.CreateAlertBody{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.CreateAlertBody{}, err
	}

	var ListAlertResp models.ListAlertsResponse
	err = json.Unmarshal(body, &ListAlertResp)
	if err != nil {
		return []models.CreateAlertBody{}, err
	}

	if !ListAlertResp.IsSuccessful {
		return []models.CreateAlertBody{}, errors.New("failed to delete view")
	}

	return ListAlertResp.Data, err
}

func DeleteAlert(client *http.Client, token string, endpoint string, teamId string, alertId []string) error {
	data := models.DeleteAlertRequest{
		AlertIds: alertId,
	}

	reqBody, err := json.Marshal(data)

	req, err := http.NewRequest("DELETE", endpoint+"api/team/"+teamId+"/alert", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "Logfire-cli")
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

	var DeleteAlertResp models.DeleteAlertResponse
	err = json.Unmarshal(body, &DeleteAlertResp)
	if err != nil {
		return err
	}

	println(DeleteAlertResp.IsSuccessful)

	if !DeleteAlertResp.IsSuccessful {
		return errors.New("failed to delete alerts")
	}

	return nil
}

func PauseAlert(client *http.Client, token string, endpoint string, teamId string, alertId []string, alertPause bool) error {
	data := models.PauseAlertRequest{
		AlertIds:   alertId,
		AlertPause: alertPause,
	}

	reqBody, err := json.Marshal(data)

	req, err := http.NewRequest("POST", endpoint+"api/team/"+teamId+"/alertpause", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "Logfire-cli")
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

	var DeleteAlertResp models.DeleteAlertResponse
	err = json.Unmarshal(body, &DeleteAlertResp)
	if err != nil {
		return err
	}

	if !DeleteAlertResp.IsSuccessful {
		if alertPause == false {
			return errors.New("failed to unpause alerts")
		} else {
			return errors.New("failed to pause alerts")
		}
	}

	return nil
}

func UpdateAlert(client *http.Client, token string, endpoint string, teamId string, name string, viewId string, numberOfRecords uint32,
	withInSeconds uint32, IntegrationsId []string, alertId string) error {

	integrationsList, err := GetAlertIntegrations(client, token, endpoint, teamId)
	if err != nil {
		return err
	}

	var integrationParsed []models.AlertIntegrationBody

	for _, integration := range integrationsList {
		for _, integrationID := range IntegrationsId {
			if integrationID == integration.ModelId {
				integrationParsed = append(integrationParsed, integration)
			}
		}
	}

	data := models.CreateAlertRequest{
		Name:            name,
		ViewId:          viewId,
		NumberOfRecords: numberOfRecords,
		WithinSeconds:   withInSeconds,
		Integrations:    integrationParsed,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", endpoint+"api/team/"+teamId+"/alert/"+alertId, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "Logfire-cli")
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

	var CreateAlertResp models.UpdateAlertResponse
	err = json.Unmarshal(body, &CreateAlertResp)
	if err != nil {
		return err
	}

	println(req.URL.Path)

	if !CreateAlertResp.IsSuccessful {
		return errors.New("failed to update alert")
	}

	return nil
}
