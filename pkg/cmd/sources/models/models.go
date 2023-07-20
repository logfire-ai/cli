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

type SourceType uint8

// Declare related constants for each SourceType starting with index 0
const (
	Kubernetes SourceType = iota + 1 // EnumIndex = 1
	AWS                              // EnumIndex = 2
	JavaScript                       // EnumIndex = 3
	Docker
	Nginx
	Dokku
	FlyDotio
	Heroku
	Ubuntu
	Vercel
	DotNET
	Apache2
	Cloudflare
	Java
	Python
	PHP
	PostgreSQL
	Redis
	Ruby
	MongoDB
	MySQL
	HTTP
	Vector
	FluentBit
	Fluentd
	Logstash
	RSyslog
	Render
	SyslogNg
	Demo
)

var EnumMap map[SourceType]string = map[SourceType]string{
	Kubernetes: "kubernetes",
	AWS:        "aws",
	JavaScript: "javascript",
	Docker:     "docker",
	Nginx:      "nginx",
	Dokku:      "dokku",
	FlyDotio:   "fly.io",
	Heroku:     "heroku",
	Ubuntu:     "ubuntu",
	Vercel:     "vercel",
	DotNET:     ".net",
	Apache2:    "apache2",
	Cloudflare: "cloudflare",
	Java:       "java",
	Python:     "python",
	PHP:        "php",
	PostgreSQL: "postgresql",
	Redis:      "redis",
	Ruby:       "ruby",
	MongoDB:    "mongodb",
	MySQL:      "mysql",
	HTTP:       "http",
	Vector:     "vector",
	FluentBit:  "fluentbit",
	Fluentd:    "fluentd",
	Logstash:   "logstash",
	RSyslog:    "rsyslog",
	Render:     "render",
	SyslogNg:   "syslog-ng",
	Demo:       "demo",
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
	"java":       13,
	"python":     14,
	"php":        15,
	"postgresql": 16,
	"redis":      17,
	"ruby":       18,
	"mongodb":    19,
	"mysql":      20,
	"http":       21,
	"vector":     22,
	"fluentbit":  23,
	"fluentd":    24,
	"logstash":   25,
	"rsyslog":    26,
	"render":     27,
	"syslog-ng":  28,
	"demo":       29,
}

func (d SourceType) String() string {

	if d < Kubernetes || d > Demo {
		return "Unknown"
	}
	return EnumMap[d]
}

func (d SourceType) EnumIndex() int {
	return int(d)
}
