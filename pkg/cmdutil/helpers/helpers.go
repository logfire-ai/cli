package helpers

import (
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/teams/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/pre_defined_prompters"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"net/http"
)

func StringNotInArray(str string, array []string) bool {
	for _, element := range array {
		if str == element {
			return false
		}
	}
	return true
}

func TeamNameToTeamId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, name string) string {
	teamsList, err := APICalls.ListTeams(client, cfg.Get().Token, cfg.Get().EndPoint)
	if err != nil {
		return ""
	}

	if len(teamsList) == 0 {
		return ""
	}

	var matchingTeams []models.Team
	for _, team := range teamsList {
		if team.Name == name {
			matchingTeams = append(matchingTeams, team)
		}
	}

	if len(matchingTeams) == 0 {
		return "" // Or return an error if appropriate
	} else if len(matchingTeams) == 1 {
		return matchingTeams[0].ID
	} else {
		fmt.Fprint(io.Out, "You seem to have multiple teams with the same name, please select from the below options")
		teamId, err := pre_defined_prompters.AskTeamId(client, cfg, io, cs, prompter)
		if err != nil {
			return ""
		}
		return teamId
	}
}
