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
)

func AskTeamId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter) (string, error) {
	teamsList, err := APICalls.ListTeams(client, cfg.Get().Token, cfg.Get().EndPoint)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s %s\n", cs.FailureIcon(), err)
		os.Exit(1)
	}

	var teamsListIdNames []string
	var idMap map[string]string = make(map[string]string)

	for _, team := range teamsList {
		lastFour := ""
		if len(team.ID) > 4 {
			lastFour = team.ID[len(team.ID)-4:]
		} else {
			lastFour = team.ID
		}
		displayName := team.Name + " - " + lastFour
		teamsListIdNames = append(teamsListIdNames, displayName)
		idMap[displayName] = team.ID
	}

	selectedTeam, err := prompter.Select("Select your Team:", "", teamsListIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	fullTeamId, ok := idMap[selectedTeam]
	if !ok {
		fmt.Fprintf(io.ErrOut, "%s Failed to map to original ID\n", cs.FailureIcon())
		return "", fmt.Errorf("Failed to map to original ID")
	}

	return fullTeamId, nil
}

func AskAlertIntegrationIds(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) ([]string, error) {
	var integrationsIdNames []string
	var integrationsList []models.AlertIntegrationBody
	var idMap map[string]string = make(map[string]string)

	integrationsList, err := APICalls.GetAlertIntegrations(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Integrations list\n", cs.FailureIcon())
		return []string{}, err
	}

	for _, integration := range integrationsList {
		lastFour := ""
		if len(integration.ModelId) > 4 {
			lastFour = integration.ModelId[len(integration.ModelId)-4:]
		} else {
			lastFour = integration.ModelId
		}
		displayName := integration.Name + " - " + lastFour
		integrationsIdNames = append(integrationsIdNames, displayName)
		idMap[displayName] = integration.ModelId
	}

	integrationsSelected, err := prompter.MultiSelect("Select integrations to be alerted. (multiple selections are allowed)", []string{}, integrationsIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return []string{}, err
	}

	var fullIntegrationIds []string
	for _, integrationSelected := range integrationsSelected {
		fullId, ok := idMap[integrationSelected]
		if !ok {
			fmt.Fprintf(io.ErrOut, "%s Failed to map to original ID\n", cs.FailureIcon())
			return []string{}, fmt.Errorf("Failed to map to original ID")
		}
		fullIntegrationIds = append(fullIntegrationIds, fullId)
	}

	return fullIntegrationIds, nil
}

func AskViewId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	ViewsList, err := APICalls.ListView(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Views list\n", cs.FailureIcon())
		return "", err
	}

	var viewsIdNames []string
	var idMap map[string]string = make(map[string]string)

	for _, view := range ViewsList {
		lastFour := ""
		if len(view.Id) > 4 {
			lastFour = view.Id[len(view.Id)-4:]
		} else {
			lastFour = view.Id
		}
		displayName := view.Name + " - " + lastFour
		viewsIdNames = append(viewsIdNames, displayName)
		idMap[displayName] = view.Id
	}

	viewSelected, err := prompter.Select("Select a View for which alert is to be created:", "", viewsIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	fullViewId, ok := idMap[viewSelected]
	if !ok {
		fmt.Fprintf(io.ErrOut, "%s Failed to map to original ID\n", cs.FailureIcon())
		return "", fmt.Errorf("Failed to map to original ID")
	}

	return fullViewId, nil
}

func AskAlertIds(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) ([]string, error) {
	var alertIdNames []string
	var alertsList []models.CreateAlertBody
	var idMap map[string]string = make(map[string]string)

	alertsList, err := APICalls.ListAlert(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return []string{}, err
	}

	for _, alert := range alertsList {
		lastFour := ""
		if len(alert.Id) > 4 {
			lastFour = alert.Id[len(alert.Id)-4:]
		} else {
			lastFour = alert.Id
		}
		displayName := alert.Name + " - " + lastFour
		alertIdNames = append(alertIdNames, displayName)
		idMap[displayName] = alert.Id
	}

	alertsSelected, err := prompter.MultiSelect("Select alerts. (multiple selections are allowed)", []string{}, alertIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return []string{}, err
	}

	var fullAlertIds []string
	for _, alertSelected := range alertsSelected {
		fullId, ok := idMap[alertSelected]
		if !ok {
			fmt.Fprintf(io.ErrOut, "%s Failed to map to original ID\n", cs.FailureIcon())
			return []string{}, fmt.Errorf("Failed to map to original ID")
		}
		fullAlertIds = append(fullAlertIds, fullId)
	}

	return fullAlertIds, nil
}

func AskAlertId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	var alertIdNames []string
	var alertsList []models.CreateAlertBody
	var idMap map[string]string = make(map[string]string)

	alertsList, err := APICalls.ListAlert(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return "", err
	}

	for _, alert := range alertsList {
		lastFour := ""
		if len(alert.Id) > 4 {
			lastFour = alert.Id[len(alert.Id)-4:]
		} else {
			lastFour = alert.Id
		}
		displayName := alert.Name + " - " + lastFour
		alertIdNames = append(alertIdNames, displayName)
		idMap[displayName] = alert.Id
	}

	alertsSelected, err := prompter.Select("Select an alert:", "", alertIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	fullAlertId, ok := idMap[alertsSelected]
	if !ok {
		fmt.Fprintf(io.ErrOut, "%s Failed to map to original ID\n", cs.FailureIcon())
		return "", fmt.Errorf("Failed to map to original ID")
	}

	return fullAlertId, nil
}

func AskIntegrationId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	var integrationsIdNames []string
	var integrationsList []IntegrationModels.IntegrationBody
	var idMap map[string]string = make(map[string]string)

	integrationsList, err := APICalls.GetIntegrationsList(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Integrations list\n", cs.FailureIcon())
		return "", err
	}

	for _, integration := range integrationsList {
		lastFour := ""
		if len(integration.Id) > 4 {
			lastFour = integration.Id[len(integration.Id)-4:]
		} else {
			lastFour = integration.Id
		}
		displayName := integration.Name + " - " + lastFour
		integrationsIdNames = append(integrationsIdNames, displayName)
		idMap[displayName] = integration.Id
	}

	integrationsSelected, err := prompter.Select("Select integrations to be alerted. (multiple selections are allowed)", "", integrationsIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	fullIntegrationId, ok := idMap[integrationsSelected]
	if !ok {
		fmt.Fprintf(io.ErrOut, "%s Failed to map to original ID\n", cs.FailureIcon())
		return "", fmt.Errorf("Failed to map to original ID")
	}

	return fullIntegrationId, nil

}

func AskSourceId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	idMap := make(map[string]string)

	var sourceIdNames []string
	var sourceList []SourceModels.Source

	sourceList, err := APICalls.GetAllSources(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return "", err
	}

	for _, source := range sourceList {
		lastFour := ""
		if len(source.ID) > 4 {
			lastFour = source.ID[len(source.ID)-4:]
		} else {
			lastFour = source.ID
		}
		sourceIdNames = append(sourceIdNames, source.Name+" - "+lastFour)
		idMap[source.Name+" - "+lastFour] = source.ID
	}

	sourceSelected, err := prompter.Select("Select a source:", "", sourceIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	fullID, ok := idMap[sourceSelected]
	if !ok {
		fmt.Fprintf(io.ErrOut, "%s Failed to map to original ID\n", cs.FailureIcon())
		return "", fmt.Errorf("Failed to map to original ID")
	}

	return fullID, nil
}

func AskSourceIds(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) ([]string, error) {
	var sourceIdNames []string
	var sourceList []SourceModels.Source
	var idMap map[string]string = make(map[string]string)

	sourceList, err := APICalls.GetAllSources(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return []string{}, err
	}

	for _, source := range sourceList {
		lastFour := ""
		if len(source.ID) > 4 {
			lastFour = source.ID[len(source.ID)-4:]
		} else {
			lastFour = source.ID
		}
		displayName := source.Name + " - " + lastFour
		sourceIdNames = append(sourceIdNames, displayName)
		idMap[displayName] = source.ID
	}

	sourcesSelected, err := prompter.MultiSelect("Select sources. (multiple selections are allowed)", []string{}, sourceIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return []string{}, err
	}

	var fullSourceIds []string
	for _, sourceSelected := range sourcesSelected {
		fullId, ok := idMap[sourceSelected]
		if !ok {
			fmt.Fprintf(io.ErrOut, "%s Failed to map to original ID\n", cs.FailureIcon())
			return []string{}, fmt.Errorf("Failed to map to original ID")
		}
		fullSourceIds = append(fullSourceIds, fullId)
	}

	return fullSourceIds, nil
}

func AskMemberId(client *http.Client, cfg config.Config, io *iostreams.IOStreams, cs *iostreams.ColorScheme, prompter prompter.Prompter, teamId string) (string, error) {
	var MemberIdNames []string
	var membersList []MemberModels.TeamMemberRes
	var idMap map[string]string = make(map[string]string)

	membersList, err := APICalls.MembersList(client, cfg.Get().Token, cfg.Get().EndPoint, teamId)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to Get Alerts list\n", cs.FailureIcon())
		return "", err
	}

	for _, member := range membersList {
		lastFour := ""
		if len(member.ProfileId) > 4 {
			lastFour = member.ProfileId[len(member.ProfileId)-4:]
		} else {
			lastFour = member.ProfileId
		}
		displayName := member.FirstName + " " + member.LastName + " - " + lastFour
		MemberIdNames = append(MemberIdNames, displayName)
		idMap[displayName] = member.ProfileId
	}

	memberSelected, err := prompter.Select("Select a member:", "", MemberIdNames)
	if err != nil {
		fmt.Fprintf(io.ErrOut, "%s Failed to read selection\n", cs.FailureIcon())
		return "", err
	}

	fullProfileId, ok := idMap[memberSelected]
	if !ok {
		fmt.Fprintf(io.ErrOut, "%s Failed to map to original ID\n", cs.FailureIcon())
		return "", fmt.Errorf("Failed to map to original ID")
	}

	return fullProfileId, nil

}
