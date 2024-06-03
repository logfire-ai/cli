package models

type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
}

type UpdateProfileResponse struct {
	IsSuccessful bool `json:"isSuccessful"`
}

type UpdateFlagRequest struct {
	TeamId string `json:"teamId"`
}

type UpdateFlagResponse struct {
	TeamId       string   `json:"teamId,omitempty"`
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
}
