package pre_defined_prompters

import (
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmd/alerts/models"
	IntegrationModels "github.com/logfire-sh/cli/pkg/cmd/integrations/models"
	SourceModels "github.com/logfire-sh/cli/pkg/cmd/sources/models"
	MemberModels "github.com/logfire-sh/cli/pkg/cmd/teams/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"net/http"
	"os"
	"strings"
)

func AskTeamId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter) (string, error) {
	teamsList, err := APICalls.ListTeams(client, cfg.Get().Token, cfg.Get().EndPoint)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s %s\n", cs.FailureIcon(), err)
		os.Exit(1)
	}

	var teamsListIdNames []string

	for _, team := range teamsList {
		teamsListIdNames = append(teamsListIdNames, team.Name+" - "+team.ID)
	}

	selectedTeam, err := prompter.Select("Select your Team:", "", teamsListIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	return strings.TrimSpace(strings.Split(selectedTeam, " - ")[1]), nil
}

func AskAlertIntegrationIds(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) ([]string, error) {
	var integrationsIdNames []string
	var integrationsList []models.AlertIntegrationBody

	integrationsList, err := APICalls.GetAlertIntegrations(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Integrations list\n", cs.FailureIcon())
		return []string{}, err
	}

	for _, integration := range integrationsList {
		integrationsIdNames = append(integrationsIdNames, integration.Name+" - "+integration.ModelId)
	}

	integrationsSelected, err := prompter.MultiSelect("Select integrations to be alerted. (multiple selections are allowed)", []string{}, integrationsIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return []string{}, err
	}

	var integrationsSelectedList []string

	for _, integrationSelected := range integrationsSelected {
		parts := strings.Split(integrationSelected, " - ")
		if len(parts) > 1 {
			// Trim any leading or trailing spaces from the right part before adding to the result slice.
			integrationsSelectedList = append(integrationsSelectedList, strings.TrimSpace(parts[1]))
		}
	}

	return integrationsSelectedList, nil

}

func AskViewId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	ViewsList, err := APICalls.ListView(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Views list\n", cs.FailureIcon())
		return "", err
	}

	var viewsIdNames []string

	for _, view := range ViewsList {
		viewsIdNames = append(viewsIdNames, view.Name+" - "+view.Id)
	}

	viewSelected, err := prompter.Select("Select a View for which alert is to be created:", "", viewsIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	return strings.TrimSpace(strings.Split(viewSelected, " - ")[1]), nil
}

func AskAlertIds(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) ([]string, error) {
	var alertIdNames []string
	var alertsList []models.CreateAlertBody

	alertsList, err := APICalls.ListAlert(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return []string{}, err
	}

	for _, integration := range alertsList {
		alertIdNames = append(alertIdNames, integration.Name+" - "+integration.Id)
	}

	alertsSelected, err := prompter.MultiSelect("Select alerts. (multiple selections are allowed)", []string{}, alertIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return []string{}, err
	}

	var alertsSelectedList []string

	for _, alertSelected := range alertsSelected {
		parts := strings.Split(alertSelected, " - ")
		if len(parts) > 1 {
			// Trim any leading or trailing spaces from the right part before adding to the result slice.
			alertsSelectedList = append(alertsSelectedList, strings.TrimSpace(parts[1]))
		}
	}

	return alertsSelectedList, nil

}

func AskAlertId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	var alertIdNames []string
	var alertsList []models.CreateAlertBody

	alertsList, err := APICalls.ListAlert(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return "", err
	}

	for _, integration := range alertsList {
		alertIdNames = append(alertIdNames, integration.Name+" - "+integration.Id)
	}

	alertsSelected, err := prompter.Select("Select an alert:", "", alertIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	return strings.TrimSpace(strings.Split(alertsSelected, " - ")[1]), nil

}

func AskIntegrationIds(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	var integrationsIdNames []string
	var integrationsList []IntegrationModels.IntegrationBody

	integrationsList, err := APICalls.GetIntegrationsList(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Integrations list\n", cs.FailureIcon())
		return "", err
	}

	for _, integration := range integrationsList {
		integrationsIdNames = append(integrationsIdNames, integration.Name+" - "+integration.Id)
	}

	integrationsSelected, err := prompter.Select("Select integrations to be alerted. (multiple selections are allowed)", "", integrationsIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	return strings.TrimSpace(strings.Split(integrationsSelected, " - ")[1]), nil
}

func AskSourceId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	var sourceIdNames []string
	var sourceList []SourceModels.Source

	sourceList, err := APICalls.GetAllSources(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return "", err
	}

	for _, source := range sourceList {
		sourceIdNames = append(sourceIdNames, source.Name+" - "+source.ID)
	}

	sourceSelected, err := prompter.Select("Select a source:", "", sourceIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	return strings.TrimSpace(strings.Split(sourceSelected, " - ")[1]), nil
}

func AskSourceIds(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) ([]string, error) {
	var sourceIdNames []string
	var sourceList []SourceModels.Source

	sourceList, err := APICalls.GetAllSources(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return []string{}, err
	}

	for _, source := range sourceList {
		sourceIdNames = append(sourceIdNames, source.Name+" - "+source.ID)
	}

	sourcesSelected, err := prompter.MultiSelect("Select sources. (multiple selections are allowed)", []string{}, sourceIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return []string{}, err
	}

	var sourceSelectedList []string

	for _, sourceSelected := range sourcesSelected {
		parts := strings.Split(sourceSelected, " - ")
		if len(parts) > 1 {
			// Trim any leading or trailing spaces from the right part before adding to the result slice.
			sourceSelectedList = append(sourceSelectedList, strings.TrimSpace(parts[1]))
		}
	}

	return sourceSelectedList, nil

}

func AskMemberId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	var MemberIdNames []string
	var membersList []MemberModels.TeamMemberRes

	membersList, err := APICalls.MembersList(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return "", err
	}

	for _, member := range membersList {
		MemberIdNames = append(MemberIdNames, member.FirstName+" "+member.LastName+" - "+member.ProfileId)
	}

	memberSelected, err := prompter.Select("Select a member:", "", MemberIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	return strings.TrimSpace(strings.Split(memberSelected, " - ")[1]), nil

}
