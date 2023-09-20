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
	Data         Source   `json:"data,omitempty"`
	Message      []string `json:"message,omitempty"`
}

type SourcesResponse struct {
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

var PlatformMap map[string]int = map[string]int{
	"kubernetes": 1,
	"aws":        2,
	"javascript": 3,
	"docker":     4,
	"nginx":      5,
	"dokku":      6,
	"fly.io":     7,
	"heroku":     8,
	"ubuntu":     9,
	"vercel":     10,
	".net":       11,
	"apache2":    12,
	"cloudflare": 13,
	"java":       14,
	"python":     15,
	"php":        16,
	"postgresql": 17,
	"redis":      18,
	"ruby":       19,
	"mongodb":    20,
	"mysql":      21,
	"http":       22,
	"vector":     23,
	"fluentbit":  24,
	"fluentd":    25,
	"logstash":   26,
	"rsyslog":    27,
	"render":     28,
	"syslog-ng":  29,
	"demo":       30,
}

type ConfigurationResponse struct {
	IsSuccessful bool        `json:"isSuccessful"`
	Message      []string    `json:"message,omitempty"`
	Data         interface{} `json:"data"`
}
