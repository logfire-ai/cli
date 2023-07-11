package models

import "time"

type Source struct {
	ID          string     `json:"id"`
	ProfileID   string     `json:"profileId"`
	TeamID      string     `json:"teamId"`
	Name        string     `json:"name"`
	SourceType  int        `json:"sourceType"`
	SourceToken string     `json:"sourceToken"`
	Platform    string     `json:"platform"`
	CreatedAt   time.Time  `json:"CreatedAt"`
	UpdatedAt   time.Time  `json:"UpdatedAt"`
	DeletedAt   *time.Time `json:"DeletedAt"`
}

type SourceResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Data         []Source `json:"data,omitempty"`
	Message      []string `json:"message,omitempty"`
}

type SourceCreateResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Data         Source   `json:"data,omitempty"`
	Message      []string `json:"message,omitempty"`
}

type SourceCreate struct {
	Name       string `json:"name"`
	SourceType int    `json:"sourceType"`
}
