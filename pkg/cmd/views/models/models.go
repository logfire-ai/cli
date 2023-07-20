package models

import (
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"time"
)

type DeleteViewResponse struct {
	IsSuccessful bool   `json:"isSuccessful"`
	Message      string `json:"message,omitempty"`
}

type DateInterval struct {
	StartDate time.Time `json:"startDate,omitempty"`
	EndDate   time.Time `json:"endDate,omitempty"`
}

type LevelObj struct {
	Id    int    `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}

type SearchObj struct {
	Condition string `json:"condition,omitempty"`
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
}

type ViewResponseBody struct {
	Name          string          `json:"name" validate:"required"`
	Description   string          `json:"description,omitempty"`
	Id            string          `json:"id,omitempty"`
	SourcesFilter []models.Source `json:"sourcesFilter,omitempty"`
	LevelFilter   *[]LevelObj     `json:"levelFilter,omitempty"`
	DateFilter    DateInterval    `json:"dateFilter,omitempty"`
	SqlFilter     string          `json:"sqlFilter,omitempty"`
	SearchFilter  []SearchObj     `json:"searchFilter,omitempty"`
	TextFilter    []string        `json:"textFilter,omitempty"`
}

type CreateViewResponse struct {
	IsSuccessful bool              `json:"isSuccessful"`
	View         *ViewResponseBody `json:"view,omitempty"`
	Message      []string          `json:"message,omitempty"`
}

type ViewResponse struct {
	IsSuccessful bool             `json:"isSuccessful"`
	Data         ViewResponseBody `json:"data"`
}

type ListViewResponse struct {
	IsSuccessful bool               `json:"isSuccessful"`
	Views        []ViewResponseBody `json:"data"`
}
