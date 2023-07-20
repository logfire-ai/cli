package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	ResetPasswordModels "github.com/logfire-sh/cli/pkg/cmd/reset_password/models"
	UpdateProfileModels "github.com/logfire-sh/cli/pkg/cmd/update_profile/models"
	"io"
	"net/http"
)

func ResetPassword(client *http.Client, token string, endpoint string, profileId string, password string) error {
	data := ResetPasswordModels.ResetPasswordRequest{
		Password: password,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", endpoint+"api/profile/"+profileId+"/set-password", bytes.NewBuffer(reqBody))
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

	var ResetPasswordResp ResetPasswordModels.ResetPasswordResponse
	err = json.Unmarshal(body, &ResetPasswordResp)
	if err != nil {
		return err
	}

	if !ResetPasswordResp.IsSuccessful {
		return errors.New("failed to change password")
	}

	return nil
}

func UpdateProfile(client *http.Client, token string, endpoint string, profileId string, firstName string, lastName string) error {
	data := UpdateProfileModels.UpdateProfileRequest{
		FirstName: firstName,
		LastName:  lastName,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", endpoint+"api/profile/"+profileId, bytes.NewBuffer(reqBody))
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

	var ResetPasswordResp UpdateProfileModels.UpdateProfileResponse
	err = json.Unmarshal(body, &ResetPasswordResp)
	if err != nil {
		return err
	}

	if !ResetPasswordResp.IsSuccessful {
		return errors.New("failed to update profile")
	}

	return nil
}
