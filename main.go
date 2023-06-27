package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"logfire/gui"
	"logfire/models"
	"logfire/sources"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SignupRequest struct {
	Email string `json:"email"`
}

type SigninRequest struct {
	AuthType   int    `json:"authType"`
	Credential string `json:"credential"`
}

type OnboardRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type SetPassword struct {
	Password string `json:"password"`
}

type UserBody struct {
	ProfileID string `json:"profileId"`
	TeamID    string `json:"teamId"`
	Onboarded bool   `json:"onboarded"`
	Email     string `json:"email"`
}

type BearerToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Exp          string `json:"exp"`
	Iat          string `json:"iat"`
}

type Response struct {
	IsSuccessful bool        `json:"isSuccessful"`
	Code         int         `json:"code"`
	Email        string      `json:"email"`
	UserBody     UserBody    `json:"userBody"`
	BearerToken  BearerToken `json:"bearerToken"`
	Message      []string    `json:"message"`
}

type SigninPasswordRequest struct {
	Email      string `json:"email"`
	AuthType   int    `json:"authType"`
	Credential string `json:"credential"`
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "logfire",
		Short: "A simple CLI application to interact with logfire.sh",
	}

	registerCmd := &cobra.Command{
		Use:   "register [config_file]",
		Short: "Register the application",
		Args:  cobra.ExactArgs(1),
		Run:   register,
	}

	signinCmd := &cobra.Command{
		Use:   "login [config_file]",
		Short: "login to the application",
		Args:  cobra.ExactArgs(1),
		Run:   signinPassword,
	}

	sourceCmd := &cobra.Command{
		Use:   "sources [list/create/delete] [config_file]",
		Short: "manage the sources",
		Args:  cobra.ExactArgs(2),
		Run:   sourceManage,
	}

	livetailCmd := &cobra.Command{
		Use:   "livetail ",
		Short: "display the livetail",
		Args:  cobra.ExactArgs(0),
		Run:   livetailShow,
	}

	rootCmd.AddCommand(registerCmd, signinCmd, sourceCmd, livetailCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func onboardProfile(response Response, host string, port int) {

	fmt.Println("Enter your First Name:")
	reader := bufio.NewReader(os.Stdin)
	firstName, _ := reader.ReadString('\n')
	lastName, _ := reader.ReadString('\n')

	onboardReq := OnboardRequest{
		FirstName: firstName,
		LastName:  lastName,
	}

	reqBody, err := json.Marshal(onboardReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	url := fmt.Sprintf("http://%s:%d/api/profile/%s/onboard", host, port, response.UserBody.ProfileID)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Sign-in successful!")
	} else {
		fmt.Printf("Sign-in failed with status code: %d\n", resp.StatusCode)
	}

}

func setPassword(response Response, host string, port int) {

	fmt.Println("Enter your Password:")
	reader := bufio.NewReader(os.Stdin)
	password, _ := reader.ReadString('\n')

	pwdReq := SetPassword{
		Password: password,
	}

	reqBody, err := json.Marshal(pwdReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	url := fmt.Sprintf("http://%s:%d/api/profile/%s/set-password", host, port, response.UserBody.ProfileID)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Sign-in successful!")
	} else {
		fmt.Printf("Sign-in failed with status code: %d\n", resp.StatusCode)
	}

}

func signin(host string, port int) *Response {

	fmt.Println("Thank You for Registering. An email has been sent to your address.\nPlease Enter the token in the email link here: ")
	reader := bufio.NewReader(os.Stdin)
	credential, _ := reader.ReadString('\n')
	fmt.Printf("Hello, %s!\n", credential)

	signinReq := SigninRequest{
		AuthType:   1,
		Credential: credential,
	}
	reqBody, err := json.Marshal(signinReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return nil
	}

	url := fmt.Sprintf("http://%s:%d/api/auth/signin", host, port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Sign-in successful!")
	} else {
		fmt.Printf("Sign-in failed with status code: %d\n", resp.StatusCode)
	}

	var response Response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return nil
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		return nil
	}

	return &response

}

func register(cmd *cobra.Command, args []string) {
	configFile := args[0]

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read configuration file: %v\n", err)
		return
	}

	host := viper.GetString("host")
	port := viper.GetInt("port")
	debug := viper.GetBool("debug")

	fmt.Printf("Registering with host: %s, port: %d, debug: %t\n", host, port, debug)

	signupReq := SignupRequest{
		Email: viper.GetString("email"),
	}
	reqBody, err := json.Marshal(signupReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	url := fmt.Sprintf("http://%s:%d/api/auth/signup", host, port)

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
		fmt.Printf("Failed to send GET request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Registration successful!")
	} else {
		fmt.Printf("Registration failed with status code: %d\n", resp.StatusCode)
	}

	fmt.Println("Thank You for Registering. An email has been sent to your adress \n Please Enter the token in the email link here:")
	reader := bufio.NewReader(os.Stdin)
	credential, _ := reader.ReadString('\n')
	fmt.Printf("Hello, %s!\n , %d", credential, len(credential))

	signinReq := SigninRequest{
		AuthType:   1,
		Credential: strings.TrimSpace(credential),
	}
	reqBody, err = json.Marshal(signinReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	url = fmt.Sprintf("http://%s:%d/api/auth/signin", host, port)
	resp, err = client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Sign-in successful!")
	} else {
		fmt.Printf("Sign-in failed with status code: %d\n", resp.StatusCode)
	}

	var response Response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		return
	}

	fmt.Println("Enter your First Name:")
	firstName, _ := reader.ReadString('\n')
	lastName, _ := reader.ReadString('\n')

	onboardReq := OnboardRequest{
		FirstName: firstName,
		LastName:  lastName,
	}

	reqBody, err = json.Marshal(onboardReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	url = fmt.Sprintf("http://%s:%d/api/profile/%s/onboard", host, port, response.UserBody.ProfileID)

	client1 := &http.Client{}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	token := response.BearerToken.AccessToken
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client1.Do(req)

	//resp, err = client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		fmt.Println("Sign-in successful!")
	} else {
		fmt.Printf("Sign-in failed with status code: %d\n", resp.StatusCode)
	}

	fmt.Println("Enter your Password:")
	password, _ := reader.ReadString('\n')

	pwdReq := SetPassword{
		Password: password,
	}

	reqBody, err = json.Marshal(pwdReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	url = fmt.Sprintf("http://%s:%d/api/profile/%s/set-password", host, port, response.UserBody.ProfileID)

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client1.Do(req)

	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		fmt.Println("Registration successful!")
	} else {
		fmt.Printf("Sign-in failed with status code: %d\n", resp.StatusCode)
	}

}

func signinPassword(cmd *cobra.Command, args []string) {

	configFile := args[0]

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read configuration file: %v\n", err)
		return
	}

	host := viper.GetString("host")
	port := viper.GetInt("port")
	debug := viper.GetBool("debug")

	fmt.Printf("Signing in with host: %s, port: %d, debug: %t\n", host, port, debug)

	fmt.Println("Enter your Password:")
	reader := bufio.NewReader(os.Stdin)
	password, _ := reader.ReadString('\n')

	signinReq := SigninPasswordRequest{
		Email:      viper.GetString("email"),
		AuthType:   2,
		Credential: strings.TrimSuffix(password, "\n"),
	}

	reqBody, err := json.Marshal(signinReq)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return
	}

	url := fmt.Sprintf("http://%s:%d/api/auth/signin", host, port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Sign-in successful!")
	} else {
		fmt.Printf("Sign-in failed with status code: %d\n", resp.StatusCode)
	}

	var response Response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		return
	}
}

func sourceManage(cmd *cobra.Command, args []string) {

	configFile := args[1]

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read configuration file: %v\n", err)
		return
	}

	host := viper.GetString("host")
	port := viper.GetInt("port")

	switch subCmd := args[0]; subCmd {
	case "list":
		fmt.Println("Enter your Token:")
		reader := bufio.NewReader(os.Stdin)
		token, _ := reader.ReadString('\n')

		fmt.Println("Enter your TeamId:")
		reader = bufio.NewReader(os.Stdin)
		teamId, _ := reader.ReadString('\n')
		teamId = strings.TrimSuffix(teamId, "\n")

		url := fmt.Sprintf("http://%s:%d/api/team/", host, port)

		url += teamId + "/source"

		resp, err := sources.GetAllSources(strings.TrimSuffix(token, "\n"), strings.TrimSuffix(teamId, "\n"), url)
		if err == nil {
			fmt.Printf("Source: %+v\n", resp)
		}
	case "create":
		fmt.Println("Enter your Token:")
		reader := bufio.NewReader(os.Stdin)
		token, _ := reader.ReadString('\n')

		fmt.Println("Enter your TeamId:")
		reader = bufio.NewReader(os.Stdin)
		teamId, _ := reader.ReadString('\n')

		fmt.Println("Enter your Source Name:")
		reader = bufio.NewReader(os.Stdin)
		name, _ := reader.ReadString('\n')

		fmt.Println("Enter your Source Type:")
		reader = bufio.NewReader(os.Stdin)
		sourceType, _ := reader.ReadString('\n')
		num, err := strconv.Atoi(sourceType)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		data := models.SourceCreate{
			Name:       name,
			SourceType: num,
		}
		url := fmt.Sprintf("http://%s:%d/api/team/", host, port)
		url += teamId + "/source"

		resp, err := sources.CreateSources(strings.TrimSuffix(token, "\n"), strings.TrimSuffix(teamId, "\n"), url, data)
		if err == nil {
			fmt.Printf("Source: %+v\n", resp)
		}
	case "delete":
		fmt.Println("Enter your Token:")
		reader := bufio.NewReader(os.Stdin)
		token, _ := reader.ReadString('\n')

		fmt.Println("Enter your TeamId:")
		reader = bufio.NewReader(os.Stdin)
		teamId, _ := reader.ReadString('\n')

		fmt.Println("Enter your Source Name:")
		reader = bufio.NewReader(os.Stdin)
		id, _ := reader.ReadString('\n')

		url := fmt.Sprintf("http://%s:%d/api/team/", host, port)
		url += teamId + "/source/" + id

		resp, err := sources.DeleteSources(strings.TrimSuffix(token, "\n"), strings.TrimSuffix(teamId, "\n"), url)
		if err == nil {
			fmt.Printf("Source: %+v\n", resp)
		}
	}

}

func livetailShow(cmd *cobra.Command, args []string) {

	configFile := args[1]

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read configuration file: %v\n", err)
		return
	}

	// host := viper.GetString("host")
	// port := viper.GetInt("port")

	fmt.Println("Enter your Token:")
	reader := bufio.NewReader(os.Stdin)
	token, _ := reader.ReadString('\n')

	fmt.Println("Enter your TeamId:")
	reader = bufio.NewReader(os.Stdin)
	teamId, _ := reader.ReadString('\n')
	teamId = strings.TrimSuffix(teamId, "\n")

	ui := gui.NewUI(token, teamId)
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

	// url := fmt.Sprintf("http://%s:%d/api/team/", host, port)

	// url += teamId + "/source"

	// err = livetail.ShowLivetail(strings.TrimSuffix(token, "\n"), strings.TrimSuffix(teamId, "\n"))
	// if err == nil {
	// 	fmt.Printf("LiveTail displayed\n")
	// }

}
