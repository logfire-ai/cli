package models

import (
	"github.com/google/uuid"
	"github.com/logfire-sh/cli/pkg/cmd/teams/constants"
)

type CreateTeamRequest struct {
	Name string  `json:"name" validate:"required"`
	Logo *string `json:"logo,omitempty"`
}

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type AllTeamResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
	Data         []Team   `json:"data,omitempty"`
}

type CreateTeamResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Data         Team     `json:"data,omitempty"`
	Message      []string `json:"message,omitempty"`
}

type AllTeamMemberResponse struct {
	IsSuccessful bool       `json:"isSuccessful"`
	Message      []string   `json:"message,omitempty"`
	Data         AllTMandTI `json:"data,omitempty"`
}

type AllTMandTI struct {
	CountTeamMembers int             `json:"countTeamMembers,omitempty"`
	CountTeamInvites int             `json:"countTeamInvites,omitempty"`
	TeamMembers      []TeamMemberRes `json:"teamMembers,omitempty"`
	TeamInvites      []TeamInviteRes `json:"teamInvites,omitempty"`
}

type TeamMemberRes struct {
	TeamMember
	// Name string `json:"name"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type RemoveMemberReq struct {
	MemberId string `json:"memberId,omitempty"`
}

type UpdateMemberReq struct {
	RemoveMemberReq
	Role int `json:"role"`
}

type UpdateMemberRes struct {
	TeamMemberRes
	IsSuccessful bool `json:"isSuccessful"`
}

type TeamMember struct {
	ProfileId string             `json:"profileId"`
	TeamId    string             `json:"teamId"`
	Role      constants.RoleType `json:"role"`
}

type TeamInviteRes struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
	//Data         []TeamInvite `json:"data,omitempty"`
}

type TeamInviteReq struct {
	Email []string `json:"email" validate:"required"`
}

type TeamInvite struct {
	ID          uuid.UUID `json:"id"`
	MagicLinkId uuid.UUID `json:"magicLinkId"`
	TeamId      uuid.UUID `json:"teamId"`
	Accepted    bool      `json:"accepted"`
	Email       string    `json:"email"`
}
