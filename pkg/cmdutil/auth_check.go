package cmdutil

import (
	"github.com/logfire-sh/cli/internal/config"
	"github.com/spf13/cobra"
)

func CheckAuth(cfg config.Config) bool {
	return cfg.HasEnvToken()
}

func IsAuthCheckEnabled(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "help", "signup", "login", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
		return false
	}

	for c := cmd; c.Parent() != nil; c = c.Parent() {
		if c.Annotations != nil && c.Annotations["skipAuthCheck"] == "true" {
			return false
		}
	}

	return true
}
