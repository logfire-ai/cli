package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/gui"
	"github.com/logfire-sh/cli/models"
	"github.com/logfire-sh/cli/pkg/cmd/auth"
	"github.com/logfire-sh/cli/sources"

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
		Use:   "logfire <command> <subcommand> [flags]",
		Short: "Logfire CLI",
		Long:  `Work seamlessly with Logfire.sh log management system from the command line.`,
		Example: heredoc.Doc(`
			$ logfire auth login
			$ logfire livetail show
		`),
	}

	rootCmd.AddGroup(&cobra.Group{
		ID:    "core",
		Title: "Core commands",
	})

	sourceCmd := &cobra.Command{
		Use:   "sources [list/create/delete] [config_file]",
		Short: "manage the sources",
		Args:  cobra.ExactArgs(2),
		Run:   sourceManage,
	}

	livetailCmd := &cobra.Command{
		Use:   "livetail ",
		Short: "display the livetail",
		Args:  cobra.ExactArgs(1),
		Run:   livetailShow,
	}

	rootCmd.AddCommand(auth.NewCmdAuth())
	rootCmd.AddCommand(sourceCmd, livetailCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
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

	configFile := args[0]

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
	token = strings.TrimSuffix(token, "\n")
	token = strings.TrimSuffix(token, "\r")

	fmt.Println("Enter your TeamId:")
	reader = bufio.NewReader(os.Stdin)
	teamId, _ := reader.ReadString('\n')
	teamId = strings.TrimSuffix(teamId, "\n")
	teamId = strings.TrimSuffix(teamId, "\r")

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
