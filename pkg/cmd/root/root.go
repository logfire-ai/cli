package root

import (
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/pkg/cmd/alerts"
	"github.com/logfire-sh/cli/pkg/cmd/integrations"
	"github.com/logfire-sh/cli/pkg/cmd/reset_password"
	"github.com/logfire-sh/cli/pkg/cmd/views"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/pkg/cmd/login"
	"github.com/logfire-sh/cli/pkg/cmd/logout"
	"github.com/logfire-sh/cli/pkg/cmd/signup"
	"github.com/logfire-sh/cli/pkg/cmd/sources"
	"github.com/logfire-sh/cli/pkg/cmd/stream"
	"github.com/logfire-sh/cli/pkg/cmd/teams"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdRoot(f *cmdutil.Factory) (*cobra.Command, error) {
	io := f.IOStreams
	cfg, err := f.Config()
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %s\n", err)
	}

	cmd := &cobra.Command{
		Use:   "logfire <command> <subcommand> [flags]",
		Short: "Logfire CLI",
		Long:  `Work seamlessly with logfire.sh log management system from the command line.`,
		Example: heredoc.Doc(`
			$ logfire login
			$ logfire stream livetail
			
		`),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// require that the user is authenticated before running most commands
			if cmdutil.IsAuthCheckEnabled(cmd) && !cmdutil.CheckAuth(cfg) {
				fmt.Fprint(io.ErrOut, authHelp())
				return errors.New("authentication required")
			}
			return nil
		},
	}

	cmd.PersistentFlags().Bool("help", false, "Show help for command")

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		rootHelpFunc(f, c, args)
	})

	cmd.SetUsageFunc(func(c *cobra.Command) error {
		return rootUsageFunc(f.IOStreams.ErrOut, c)
	})

	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	cmd.AddGroup(&cobra.Group{
		ID:    "core",
		Title: "Core commands",
	})

	cmd.AddCommand(signup.NewSignupCmd(f))
	cmd.AddCommand(login.NewLoginCmd(f))
	cmd.AddCommand(reset_password.NewResetPasswordCmd(f))
	cmd.AddCommand(logout.NewLogoutCmd(f))
	cmd.AddCommand(sources.NewCmdSource(f))
	cmd.AddCommand(teams.NewCmdTeam(f))
	cmd.AddCommand(stream.NewCmdStream(f))
	cmd.AddCommand(views.NewCmdViews(f))
	cmd.AddCommand(alerts.NewCmdAlerts(f))
	cmd.AddCommand(integrations.NewCmdIntegrations(f))
	return cmd, nil
}
