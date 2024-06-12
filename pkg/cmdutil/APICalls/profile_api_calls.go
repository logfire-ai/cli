package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/logfire-sh/cli/internal/config"
	LoginModels "github.com/logfire-sh/cli/pkg/cmd/login/models"
	ResetPasswordModels "github.com/logfire-sh/cli/pkg/cmd/reset_password/models"
	SignupModels "github.com/logfire-sh/cli/pkg/cmd/signup/models"
	UpdateProfileModels "github.com/logfire-sh/cli/pkg/cmd/update_profile/models"

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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

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

func ResetPassword(token string, endpoint string, profileId string, password string) error {
	client := http.Client{}

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
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
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

func SetPassword(token string, endpoint string, profileId string, password string) error {
	client := http.Client{}

	data := ResetPasswordModels.ResetPasswordRequest{
		Password: password,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint+"api/profile/"+profileId+"/set-password", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
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
		return errors.New(ResetPasswordResp.Message[0])
	}

	return nil
}

func UpdateProfile(client *http.Client, token string, endpoint string, profileId string, firstName string, lastName string, role string) error {
	data := UpdateProfileModels.UpdateProfileRequest{
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
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
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
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

func SignupFlow(email, endpoint string) (string, error) {
	signupReq := SignupModels.SignupRequest{
		Email: email,
	}

	reqBody, err := json.Marshal(signupReq)
	if err != nil {
		return "", err
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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var response SignupModels.SignupResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return "", err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		return "", err
	}

	if !response.IsSuccessful {
		return "", errors.New(response.Message[0])
	} else if len(response.Message) > 0 {
		return response.Message[0], nil
	} else {
		return "", nil
	}
}

func OnboardingFlow(profileID, authToken, endpoint, firstName, lastName, role string) error {
	var response LoginModels.Response

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

	client := http.Client{}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("User-Agent", "Logfire-cli")
	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		return err
	}

	if !response.IsSuccessful {
		return errors.New(response.Message[0])
	}

	return nil
}

func TokenSignIn(cfg config.Config, token, endpoint string) (error, bool) {
	var response LoginModels.Response

	signinReq := LoginModels.SigninRequest{
		AuthType:   1,
		Credential: strings.TrimSpace(token),
	}

	client := &http.Client{}

	reqBody, err := json.Marshal(signinReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return err, false
	}

	url := endpoint + "api/auth/signin"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err, false
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return err, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return err, false
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		return err, false
	}

	if !response.IsSuccessful {
		return errors.New(response.Message[0]), false
	}

	err = cfg.UpdateConfig(&response.UserBody.Email, &response.UserBody.Role, &response.BearerToken.AccessToken, &response.UserBody.ProfileID,
		&response.BearerToken.RefreshToken, nil, &response.UserBody.AccountID, nil, nil, nil, nil)
	if err != nil {
		fmt.Printf("Failed to update config: %v\n", err)
		return err, false
	}

	if response.UserBody.Onboarded {
		return nil, true
	}

	return nil, false
}

func UpdateFlag(cfg config.Config, profileID, teamId, endpoint string) error {
	var response UpdateProfileModels.UpdateFlagResponse

	updateFlagReq := UpdateProfileModels.UpdateFlagRequest{
		TeamId: teamId,
	}

	client := &http.Client{}

	reqBody, err := json.Marshal(updateFlagReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return err
	}

	url := endpoint + "/" + profileID + "/update-flags"

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("\nError: Connection failed (Server down or no internet)\n")
			os.Exit(1)
		}

		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		return err
	}

	if !response.IsSuccessful {
		return errors.New(response.Message[0])
	}

	// err = cfg.UpdateConfig(&response.UserBody.Email, &response.UserBody.Role, &response.BearerToken.AccessToken, &response.UserBody.ProfileID,
	// 	&response.BearerToken.RefreshToken, nil, &response.UserBody.AccountID, nil, nil, nil, nil)
	// if err != nil {
	// 	fmt.Printf("Failed to update config: %v\n", err)
	// 	return err
	// }

	return nil
}
