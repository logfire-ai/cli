package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	sourceModels "github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmd/views/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/filters"
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

func ListView(client *http.Client, token string, endpoint string, teamId string) ([]models.ViewResponseBody, error) {
	req, err := http.NewRequest("GET", endpoint+"api/team/"+teamId+"/view", nil)
	if err != nil {
		return []models.ViewResponseBody{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []models.ViewResponseBody{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.ViewResponseBody{}, err
	}

	var ListViewResp models.ListViewResponse
	err = json.Unmarshal(body, &ListViewResp)
	if err != nil {
		return []models.ViewResponseBody{}, err
	}

	if !ListViewResp.IsSuccessful {
		return []models.ViewResponseBody{}, errors.New("failed to delete view")
	}

	return ListViewResp.Views, err
}

func GetView(client *http.Client, token string, endpoint string, teamId string, viewId string) (models.ViewResponseBody, error) {
	req, err := http.NewRequest("GET", endpoint+"api/team/"+teamId+"/view/"+viewId, nil)
	if err != nil {
		return models.ViewResponseBody{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.ViewResponseBody{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ViewResponseBody{}, err
	}

	var ListViewResp models.ViewResponse
	err = json.Unmarshal(body, &ListViewResp)
	if err != nil {
		return models.ViewResponseBody{}, err
	}

	if !ListViewResp.IsSuccessful {
		return models.ViewResponseBody{}, errors.New("failed to get view")
	}

	return ListViewResp.Data, err
}

func CreateView(client *http.Client, token string, endpoint string, teamId string, sourceFilter []sourceModels.Source, searchFilter []string, fieldName, fieldValue,
	fieldCondition, startDate, endDate, viewName string) error {

	fieldFilter := []models.SearchObj{{
		Key:       fieldName,
		Value:     fieldValue,
		Condition: fieldCondition,
	}}

	dateFilter := models.DateInterval{
		StartDate: filters.ShortDateTimeToGoDate(startDate),
		EndDate:   filters.ShortDateTimeToGoDate(endDate),
	}

	data := models.ViewResponseBody{
		SourcesFilter: sourceFilter,
		TextFilter:    searchFilter,
		SearchFilter:  fieldFilter,
		DateFilter:    dateFilter,
		Name:          viewName,
	}

	reqBody, err := json.Marshal(data)

	req, err := http.NewRequest("POST", endpoint+"api/team/"+teamId+"/view", bytes.NewBuffer(reqBody))
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

	var CreateViewResp models.CreateViewResponse
	err = json.Unmarshal(body, &CreateViewResp)
	if err != nil {
		return err
	}

	return nil
}
