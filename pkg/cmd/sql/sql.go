package sql

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/pkg/cmd/sql/models"
	sql "github.com/logfire-sh/cli/pkg/cmd/sql/sql_recommend"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSql(f *cmdutil.Factory) *cobra.Command {
	opts := &models.SQLQueryOptions{
		IO: f.IOStreams,

		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "sql",
		Short: "Run a sql query",
		Long:  "Run a sql query",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire sql

			# start argument setup
			$ logfire sql --team-id <team-id> --query <query>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			cmdutil.SqlQueryRun(opts, f)
		},
	}
	cmd.Flags().StringVarP(&opts.TeamId, "team-id", "t", "", "Team id to be queried.")
	cmd.Flags().StringVarP(&opts.SQLQuery, "sql-query", "q", "", "SQL Query.")

	cmd.AddCommand(sql.NewSQLRecommendCmd(f))
	return cmd
}
