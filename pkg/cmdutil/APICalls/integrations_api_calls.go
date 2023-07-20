package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	AlertModels "github.com/logfire-sh/cli/pkg/cmd/alerts/models"
	IntegrationModels "github.com/logfire-sh/cli/pkg/cmd/integrations/models"
	"io"
	"net/http"
)

func GetAlertIntegrations(client *http.Client, token string, endpoint string, teamId string) ([]AlertModels.AlertIntegrationBody, error) {
	req, err := http.NewRequest("GET", endpoint+"api/team/"+teamId+"/alertintegrations", nil)
	if err != nil {
		return []AlertModels.AlertIntegrationBody{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []AlertModels.AlertIntegrationBody{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []AlertModels.AlertIntegrationBody{}, err
	}

	var ListAlertIntegrationsResp AlertModels.ListAlertIntegrationsResponse
	err = json.Unmarshal(body, &ListAlertIntegrationsResp)
	if err != nil {
		return []AlertModels.AlertIntegrationBody{}, err
	}

	if !ListAlertIntegrationsResp.IsSuccessful {
		return []AlertModels.AlertIntegrationBody{}, errors.New("failed to get alert integrations")
	}

	return ListAlertIntegrationsResp.Data, err
}

func CreateIntegration(client *http.Client, token string, endpoint string, teamId string, name, description, Id, integrationType string) error {

	data := IntegrationModels.CreateIntegrationRequest{
		Name:            name,
		IntegrationType: IntegrationModels.IntegrationMap[integrationType],
		Description:     description,
		Id:              Id,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint+"api/team/"+teamId+"/integration", bytes.NewBuffer(reqBody))
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

	println(string(body))

	var CreateIntegrationResp IntegrationModels.CreateIntegrationResponse
	err = json.Unmarshal(body, &CreateIntegrationResp)
	if err != nil {
		return err
	}

	if !CreateIntegrationResp.IsSuccessful {
		return errors.New("failed to create integration")
	}

	return nil
}

func GetIntegrationsList(client *http.Client, token string, endpoint string, teamId string) ([]IntegrationModels.IntegrationBody, error) {
	req, err := http.NewRequest("GET", endpoint+"api/team/"+teamId+"/integration", nil)
	if err != nil {
		return []IntegrationModels.IntegrationBody{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []IntegrationModels.IntegrationBody{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []IntegrationModels.IntegrationBody{}, err
	}

	var ListAlertIntegrationsResp IntegrationModels.ListIntegrationResponse
	err = json.Unmarshal(body, &ListAlertIntegrationsResp)
	if err != nil {
		return []IntegrationModels.IntegrationBody{}, err
	}

	if !ListAlertIntegrationsResp.IsSuccessful {
		return []IntegrationModels.IntegrationBody{}, errors.New("failed to get integrations")
	}

	return ListAlertIntegrationsResp.Data, err
}

func DeleteIntegration(client *http.Client, token string, endpoint string, teamId string, integrationId string) error {

	req, err := http.NewRequest("DELETE", endpoint+"api/team/"+teamId+"/integration/"+integrationId, nil)
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

	var DeleteAlertResp IntegrationModels.DeleteIntegrationResponse
	err = json.Unmarshal(body, &DeleteAlertResp)
	if err != nil {
		return err
	}

	if !DeleteAlertResp.IsSuccessful {
		return errors.New("failed to delete integration")
	}

	return nil
}

func UpdateIntegration(client *http.Client, token string, endpoint string, teamId string, integrationId, name, description, Id string) error {

	data := IntegrationModels.UpdateIntegrationRequest{
		Name:        name,
		Description: description,
		Id:          Id,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", endpoint+"api/team/"+teamId+"/integration/"+integrationId, bytes.NewBuffer(reqBody))
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

	println(string(body))

	var CreateIntegrationResp IntegrationModels.UpdateIntegrationResponse
	err = json.Unmarshal(body, &CreateIntegrationResp)
	if err != nil {
		return err
	}

	if !CreateIntegrationResp.IsSuccessful {
		return errors.New("failed to update integration")
	}

	return nil
}
