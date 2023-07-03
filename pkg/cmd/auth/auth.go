package auth

import (
	"github.com/logfire-sh/cli/pkg/cmd/auth/login"
	"github.com/logfire-sh/cli/pkg/cmd/auth/signup"
	"github.com/spf13/cobra"
)

func NewCmdAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth <command>",
		Short:   "Authenticate logfire.sh",
		GroupID: "core",
	}

	disableAuthCheck(cmd)
	cmd.AddCommand(login.NewLoginCmd())
	cmd.AddCommand(signup.NewSignupCmd())

	return cmd
}

func disableAuthCheck(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}

	cmd.Annotations["skipAuthCheck"] = "true"
}
