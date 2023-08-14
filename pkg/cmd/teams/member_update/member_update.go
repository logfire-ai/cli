package member_update

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

var RoleOptions = map[string]int{
	"member": 1,
	"admin":  2,
}

type MemberUpdateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	MemberId    string
	Role        string
}

func NewMemberUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &MemberUpdateOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "update-member",
		Short: "update-member role of a team",
		Long: heredoc.Docf(`
			update-member role of a team.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams update-member

			# start argument setup
			$ logfire teams update-member --teamid <team-id> --memberid <member-id> --role <admin|member>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			UpdateMemberRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "teamid", "", "Team id for which member is to be updated.")
	cmd.Flags().StringVar(&opts.MemberId, "memberid", "", "Member id of the member to be updated")
	cmd.Flags().StringVar(&opts.Role, "role", "", "role to which the member has to be updated")
	return cmd
}

func UpdateMemberRun(opts *MemberUpdateOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive && opts.TeamId == "" {
		opts.TeamId, _ = pre_defined_prompters.AskTeamId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter)

		opts.MemberId, _ = pre_defined_prompters.AskMemberId(opts.HttpClient(), cfg, opts.IO, cs, opts.Prompter, opts.TeamId)

		opts.Role, _ = opts.Prompter.Select("Select a new role:", "", []string{"admin", "member"})
	} else {
		if opts.TeamId == "" {
			fmt.Fprint(opts.IO.ErrOut, "team-id is required.")
			os.Exit(0)
		}

		if opts.MemberId == "" {
			fmt.Fprint(opts.IO.ErrOut, "member-id is required.")
			os.Exit(0)
		}

		if opts.Role == "" {
			fmt.Fprint(opts.IO.ErrOut, "role is required.")
			os.Exit(0)
		}

		if (opts.Role != "admin") && (opts.Role != "member") {
			fmt.Fprint(opts.IO.ErrOut, "role is not valid.")
			os.Exit(0)
		}
	}

	roleInt := RoleOptions[opts.Role]

	err = APICalls.UpdateMember(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, opts.TeamId, opts.MemberId, roleInt)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Member role updated successfully!\n", cs.SuccessIcon())
	}
}
