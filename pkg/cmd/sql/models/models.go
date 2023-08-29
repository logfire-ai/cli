package models

import (
	"net/http"

	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/iostreams"
)

type SQLFieldsBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SQLResponse struct {
	Records []map[string]interface{} `json:"records"`
	Fields  []SQLFieldsBody          `json:"fields"`
}

type ResponseItem struct {
	CaptionTitle       string `json:"caption_title"`
	CaptionDescription string `json:"caption_description"`
	SQLStatement       string `json:"sql_statement"`
}

type RecommendResponse struct {
	Data []ResponseItem `json:"data"`
}

type SQLQueryOptions struct {
	IO *iostreams.IOStreams

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool
	TeamId      string
	SQLQuery    string
	Role        string
}
