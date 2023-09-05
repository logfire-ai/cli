package main

import (
	"fmt"
	"os"

	surveyCore "github.com/AlecAivazis/survey/v2/core"
	"github.com/logfire-sh/cli/pkg/cmd/factory"
	"github.com/logfire-sh/cli/pkg/cmd/root"
	"github.com/mgutz/ansi"

	"github.com/spf13/cobra"
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
	cmdFactory := factory.New()

	cmdCh := make(chan bool)

	stderr := cmdFactory.IOStreams.ErrOut

	rootCmd, err := root.NewCmdRoot(cmdFactory, cmdCh)
	if err != nil {
		fmt.Fprintf(stderr, "failed to create root command: %s\n", err)
		return
	}

	if !cmdFactory.IOStreams.ColorEnabled() {
		surveyCore.DisableColor = true
		ansi.DisableColors(true)
	} else {
		// override survey's poor choice of color
		surveyCore.TemplateFuncsWithColor["color"] = func(style string) string {
			switch style {
			case "white":
				return ansi.ColorCode("default")
			default:
				return ansi.ColorCode(style)
			}
		}
	}

	expandedArgs := []string{}
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	// translate `gh help <command>` to `gh <command> --help` for extensions.
	if len(expandedArgs) >= 2 && expandedArgs[0] == "help" && isExtensionCommand(rootCmd, expandedArgs[1:]) {
		expandedArgs = expandedArgs[1:]
		expandedArgs = append(expandedArgs, "--help")
	}

	rootCmd.SetArgs(expandedArgs)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(stderr, "failed to run application: %s\n", err)
	}

	// rootCmd := &cobra.Command{
	// 	Use:   "logfire <command> <subcommand> [flags]",
	// 	Short: "Logfire CLI",
	// 	Long:  `Work seamlessly with Logfire.sh log management system from the command line.`,
	// 	Example: heredoc.Doc(`
	// 		$ logfire auth login
	// 		$ logfire livetail show
	// 	`),
	// 	// PersistentPreRun: func(cmd *cobra.Command, args []string) {
	// 	// 	// require that the user is authenticated before running most commands
	// 	// 	if cmdutil.IsAuthCheckEnabled(cmd) && !cmdutil.CheckAuth(cfg) {
	// 	// 		fmt.Fprint(io.ErrOut, authHelp())
	// 	// 	}
	// 	// },
	// }

	// rootCmd.AddGroup(&cobra.Group{
	// 	ID:    "core",
	// 	Title: "Core commands",
	// })

	// sourceCmd := &cobra.Command{
	// 	Use:   "sources [list/create/delete] [config_file]",
	// 	Short: "manage the sources",
	// 	Args:  cobra.ExactArgs(2),
	// 	Run:   sourceManage,
	// }

	// livetailCmd := &cobra.Command{
	// 	Use:   "livetail ",
	// 	Short: "display the livetail",
	// 	Args:  cobra.ExactArgs(1),
	// 	Run:   livetailShow,
	// }

	// rootCmd.AddCommand(auth.NewCmdAuth())
	// rootCmd.AddCommand(sourceCmd, livetailCmd)

	// if err := rootCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// }

}

// isExtensionCommand returns true if args resolve to an extension command.
func isExtensionCommand(rootCmd *cobra.Command, args []string) bool {
	c, _, err := rootCmd.Find(args)
	return err == nil && c != nil && c.GroupID == "extension"
}
