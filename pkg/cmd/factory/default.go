package factory

import (
	"net/http"
	"time"

	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
)

func New() *cmdutil.Factory {
	f := &cmdutil.Factory{
		Config: configFunc(), // No factory dependencies
	}

	f.IOStreams = ioStreams()       // No dependencies
	f.HttpClient = httpClientFunc() // No dependencies
	f.Prompter = newPrompter(f)     // Depends on IOStreams

	return f
}

func configFunc() func() (config.Config, error) {
	var cachedConfig config.Config
	var configError error
	return func() (config.Config, error) {
		if cachedConfig != nil || configError != nil {
			return cachedConfig, configError
		}
		cachedConfig, configError = config.NewConfig()
		return cachedConfig, configError
	}
}

func ioStreams() *iostreams.IOStreams {
	io := iostreams.System()
	return io
}

func httpClientFunc() func() *http.Client {
	return func() *http.Client {
		transport := http.Transport{
			IdleConnTimeout:   30 * time.Second,
			MaxIdleConns:      100,
			MaxConnsPerHost:   0,
			DisableKeepAlives: false,
		}

		client := http.Client{
			Transport: &transport,
			Timeout:   10 * time.Second,
		}
		return &client
	}
}

func newPrompter(f *cmdutil.Factory) prompter.Prompter {
	io := f.IOStreams
	return prompter.New(io.In, io.Out, io.ErrOut)
}
