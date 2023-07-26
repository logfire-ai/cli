package models

var IntegrationMap map[string]int = map[string]int{
	"member":  1,
	"email":   2,
	"slack":   3,
	"webhook": 4,
}

type CreateIntegrationRequest struct {
	Name            string `json:"name" validate:"required"`
	IntegrationType int    `json:"type" validate:"required"`
	Description     string `json:"description,omitempty"`
	Id              string `json:"email,omitempty"`
}

type CreateIntegrationResponse struct {
	IsSuccessful bool     `json:"is_successful" validate:"required"`
	Message      []string `json:"message,omitempty"`
}

type IntegrationBody struct {
	Name        string `json:"name" validate:"required"`
	Type        int    `json:"type" validate:"required"`
	Description string `json:"description,omitempty"`
	Email       string `json:"emailAddress,omitempty"`
	Id          string `json:"id" validate:"required"`
	TeamId      string `json:"teamId" validate:"required"`
}
type ListIntegrationResponse struct {
	IsSuccessful bool              `json:"isSuccessful" validate:"required"`
	Data         []IntegrationBody `json:"data"`
	Message      []string          `json:"message,omitempty"`
}

type DeleteIntegrationResponse struct {
	IsSuccessful bool   `json:"isSuccessful" validate:"required"`
	Message      string `json:"message,omitempty"`
}

type UpdateIntegrationRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateIntegrationResponse struct {
	IsSuccessful bool            `json:"isSuccessful" validate:"required"`
	Data         IntegrationBody `json:"data"`
	Message      []string        `json:"message,omitempty"`
}
