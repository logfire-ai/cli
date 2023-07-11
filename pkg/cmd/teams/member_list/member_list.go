package member_list

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/teams/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type MemberListOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
}

func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &MemberListOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "list-members",
		Args:  cobra.ExactArgs(0),
		Short: "List team members",
		Long: heredoc.Docf(`
			List all the members of a team.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire teams list members
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			if !opts.Interactive && opts.TeamId == "" {
				fmt.Fprint(opts.IO.ErrOut, "team-id is required.")
			}

			listMembersRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TeamId, "team-id", "", "Team id for which members are to be fetched.")
	return cmd
}

func listMembersRun(opts *MemberListOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	members, err := MembersList(opts.HttpClient(), cfg.Get().Token, opts.TeamId)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	}

	fmt.Fprintf(opts.IO.Out, "%s Team members retrieved successfully!\n", cs.SuccessIcon())
	for _, v := range members.TeamMembers {
		fmt.Fprintf(opts.IO.Out, "%s %s %s %s %s\n", cs.IntermediateIcon(), *v.FirstName, *v.LastName, v.ProfileId, v.Role)
	}
}

func MembersList(client *http.Client, token, teamId string) (models.AllTMandTI, error) {
	url := "https://api.logfire.sh/api/team/" + teamId + "/members"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.AllTMandTI{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return models.AllTMandTI{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.AllTMandTI{}, err
	}

	var response models.AllTeamMemberResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.AllTMandTI{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return models.AllTMandTI{}, err
	}

	if !response.IsSuccessful {
		return models.AllTMandTI{}, errors.New("api error")
	}

	return response.Data, nil
}
