package APICalls

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/logfire-sh/cli/pkg/cmd/teams/models"
	"io"
	"io/ioutil"
	"net/http"
)

func InviteMembers(client *http.Client, token string, endpoint string, teamId string, email []string) error {
	data := models.TeamInviteReq{
		Email: email,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint+"api/team/"+teamId+"/invites", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var InviteMemberResp models.TeamInviteRes
	err = json.Unmarshal(body, &InviteMemberResp)
	if err != nil {
		return err
	}

	if !InviteMemberResp.IsSuccessful {
		return errors.New("failed to update team")
	}

	return nil
}

func RemoveMember(client *http.Client, token string, endpoint string, teamId string, memberId string) error {
	data := models.RemoveMemberReq{
		MemberId: memberId,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", endpoint+"api/team/"+teamId+"/members", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var InviteMemberResp models.TeamInviteRes
	err = json.Unmarshal(body, &InviteMemberResp)
	if err != nil {
		return err
	}

	if !InviteMemberResp.IsSuccessful {
		return errors.New("failed to remove member from team")
	}

	return nil
}

func UpdateMember(client *http.Client, token string, endpoint string, teamId string, memberId string, role int) error {
	data := models.UpdateMemberReq{
		RemoveMemberReq: models.RemoveMemberReq{MemberId: memberId},
		Role:            role,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", endpoint+"api/team/"+teamId+"/members", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "Logfire-cli")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var InviteMemberResp models.TeamInviteRes
	err = json.Unmarshal(body, &InviteMemberResp)
	if err != nil {
		return err
	}

	if !InviteMemberResp.IsSuccessful {
		return errors.New("failed to update member role")
	}

	return nil
}

func MembersList(client *http.Client, token, endpoint string, teamId string) ([]models.TeamMemberRes, error) {
	url := endpoint + "api/team/" + teamId + "/members"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []models.TeamMemberRes{}, err
	}
	req.Header.Add("User-Agent", "Logfire-cli")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return []models.TeamMemberRes{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []models.TeamMemberRes{}, err
	}

	var response models.AllTeamMemberResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []models.TeamMemberRes{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return []models.TeamMemberRes{}, err
	}

	if !response.IsSuccessful {
		return []models.TeamMemberRes{}, errors.New("api error")
	}

	return response.Data.TeamMembers, nil
}
