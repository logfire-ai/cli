package cmdutil

import (
	"net/http"

	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/iostreams"
)

type Factory struct {
	IOStreams *iostreams.IOStreams
	Prompter  prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)
}
