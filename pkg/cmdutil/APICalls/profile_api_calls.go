package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/internal/prompter"
	LoginModels "github.com/logfire-sh/cli/pkg/cmd/login/models"
	ResetPasswordModels "github.com/logfire-sh/cli/pkg/cmd/reset_password/models"
	SignupModels "github.com/logfire-sh/cli/pkg/cmd/signup/models"
	UpdateProfileModels "github.com/logfire-sh/cli/pkg/cmd/update_profile/models"

	"github.com/logfire-sh/cli/pkg/iostreams"
	"io"
	"net/http"
	"time"
)

func SendMagicLink(endpoint, email string) error {
	signupReq := SignupModels.SignupRequest{
		Email: email,
	}

	reqBody, err := json.Marshal(signupReq)
	if err != nil {
		return err
	}

	url := endpoint + "api/auth/magiclink"

	transport := http.Transport{
		IdleConnTimeout:   30 * time.Second,
		MaxIdleConns:      100,
		MaxConnsPerHost:   0,
		DisableKeepAlives: false,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response LoginModels.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.New(response.Message[0])
	}

	return nil
}

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

func SignupFlow(email string, endpoint string) error {
	signupReq := SignupModels.SignupRequest{
		Email: email,
	}

	reqBody, err := json.Marshal(signupReq)
	if err != nil {
		return err
	}

	url := endpoint + "api/auth/signup"

	transport := http.Transport{
		IdleConnTimeout:   30 * time.Second,
		MaxIdleConns:      100,
		MaxConnsPerHost:   0,
		DisableKeepAlives: false,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return err
	}

	return nil
}

func OnboardingFlow(IO *iostreams.IOStreams, prompt prompter.Prompter, profileID, authToken string, endpoint, firstName, lastName string) error {
	onboardReq := SignupModels.OnboardRequest{
		FirstName: firstName,
		LastName:  lastName,
	}

	reqBody, err := json.Marshal(onboardReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return err
	}

	url := endpoint + "api/profile/" + profileID + "/onboard"

	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("unable to process the request")
	}

	return nil
}
