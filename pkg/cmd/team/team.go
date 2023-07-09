package team

import (
	"github.com/logfire-sh/cli/pkg/cmd/team/teamlist"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewTeamCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "teams <command>",
		Short:   "Team details",
		GroupID: "core",
	}

	cmd.AddCommand(teamlist.NewListCmd(f))

	return cmd
}

// func TeamsFlow(client *http.Client, io *iostreams.IOStreams, cfg config.Config, prmpt prompter.Prompter) error {
// 	cs := io.ColorScheme()

// 	url := "https://api.logfire.sh/api/team"
// 	resp, err := client.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return err
// 	}

// 	var response AllTeamResponse
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	err = json.Unmarshal(body, &response)
// 	if err != nil {
// 		return err
// 	}

// 	teams := make(map[string]Team)

// 	var options []string
// 	for _, v := range response.Data {
// 		options = append(options, v.Name)
// 		teams[v.Name] = v
// 	}

// 	result, err := prmpt.Select(
// 		"Please select the team?",
// 		options[0],
// 		options)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(io.Out, "\n%s You have selected %s\n", cs.SuccessIcon(), options[result])
// 	return err
// }

// func TeamsList(client *http.Client, io *iostreams.IOStreams, cfg config.Config, prmpt prompter.Prompter) error {
// 	url := "https://api.logfire.sh/api/team"

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return err
// 	}

// 	req.Header.Set("Authorization", "Bearer "+cfg.Get().Token)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return err
// 	}

// 	var response AllTeamResponse
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	err = json.Unmarshal(body, &response)
// 	if err != nil {
// 		return err
// 	}

// 	teams := response.Data
// 	for _, v := range teams {
// 		fmt.Fprintf(io.Out, "%s \n", v)
// 	}

// 	return nil
// }
