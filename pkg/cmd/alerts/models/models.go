package models

type DeleteAlertRequest struct {
	AlertIds []string `json:"alertids"`
}

type DeleteAlertResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
}

type AlertIntegrationBody struct {
	ModelId     string `json:"modelId" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Type        int    `json:"type" validate:"required"`
	Identifier  string `json:"identifier" validate:"required"`
	Description string `json:"description,omitempty"`
}

type ListAlertIntegrationsResponse struct {
	IsSuccessful bool                   `json:"isSuccessful"`
	Message      []string               `json:"message,omitempty"`
	Data         []AlertIntegrationBody `json:"data,omitempty"`
}

type CreateAlertRequest struct {
	Name                    string                 `json:"name" validate:"required"`
	Description             string                 `json:"description,omitempty"`
	ViewId                  string                 `json:"viewId"`
	AlertWhenHasMoreRecords bool                   `json:"alertWhenHasMoreRecords"`
	NumberOfRecords         uint32                 `json:"numberOfRecords"`
	WithinSeconds           uint32                 `json:"withinSeconds"`
	AlertSeverity           string                 `json:"alertSeverity"`
	Integrations            []AlertIntegrationBody `json:"integrations"`
	AlertPaused             bool                   `json:"alertPaused"`
	AlertLabels             *[]string              `json:"alertLabels,omitempty"`
	Runbook                 string                 `json:"runbook,omitempty"`
}

type CreateAlertBody struct {
	Name                    string                 `json:"name" validate:"required"`
	Description             string                 `json:"description,omitempty"`
	ViewId                  string                 `json:"viewId"`
	AlertWhenHasMoreRecords bool                   `json:"alertWhenHasMoreRecords"`
	NumberOfRecords         uint32                 `json:"numberOfRecords"`
	WithinSeconds           uint32                 `json:"withinSeconds"`
	AlertSeverity           string                 `json:"alertSeverity"`
	Integrations            []AlertIntegrationBody `json:"integrations"`
	AlertPaused             bool                   `json:"alertPaused"`
	AlertLabels             *[]string              `json:"alertLabels,omitempty"`
	Runbook                 string                 `json:"runbook,omitempty"`
	Id                      string                 `json:"id" validate:"required"`
}

type CreateAlertResponse struct {
	IsSuccessful bool            `json:"isSuccessful"`
	Data         CreateAlertBody `json:"data,omitempty"`
}

type UpdateAlertResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
}

type ListAlertsResponse struct {
	IsSuccessful bool              `json:"isSuccessful"`
	Data         []CreateAlertBody `json:"data,omitempty"`
}

type PauseAlertRequest struct {
	AlertIds   []string `json:"alertIds"`
	AlertPause bool     `json:"alertPaused"`
}
