package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"io"
	"net/http"
)

func UpdateSource(client *http.Client, token, endpoint string, teamid, sourceid, sourcename string) (models.Source, error) {
	data := models.SourceCreate{
		Name: sourcename,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return models.Source{}, err
	}

	req, err := http.NewRequest("PUT", endpoint+"api/team/"+teamid+"/source/"+sourceid, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.Source{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.Source{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Source{}, err
	}

	var sourceUpdateResp models.SourceCreateResponse
	err = json.Unmarshal(body, &sourceUpdateResp)
	if err != nil {
		return models.Source{}, err
	}

	if !sourceUpdateResp.IsSuccessful {
		return sourceUpdateResp.Data, errors.New("source update failed")
	}

	return sourceUpdateResp.Data, nil
}
