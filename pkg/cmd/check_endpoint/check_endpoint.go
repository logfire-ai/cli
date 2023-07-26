package check_endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"io"
	"net/http"
)

type CheckEndpointOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Endpoint bool
}

type CheckResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
}

func NewCheckEndpointCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CheckEndpointOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "check-urls",
		Short: "Check all end-points status",
		Long: heredoc.Docf(`
			Check all end-points status
		`, "`"),
		Example: heredoc.Doc(`
			$ logfire check-urls
		`),
		Run: func(cmd *cobra.Command, args []string) {
			CheckEndpointRun(opts)
		},
		GroupID: "core",
	}

	cmd.Flags().BoolVarP(&opts.Endpoint, "staging", "s", false, "To check for Staging Endpoint (TRUE|FALSE = default).")

	return cmd
}

func CheckEndpointRun(opts *CheckEndpointOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
		return
	}

	if opts.Endpoint {
		err := cfg.UpdateEndpoint("https://api-stg.logfire.ai/")
		if err != nil {
			return
		}
	}

	auth := CheckAuth(opts.HttpClient(), cfg.Get().EndPoint)
	if auth == "Hello From Auth!!!" {
		fmt.Fprintf(opts.IO.Out, "%s Auth endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Auth endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	profile := CheckProfile(opts.HttpClient(), cfg.Get().EndPoint)
	if !profile.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Profile endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Profile endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	source := CheckSource(opts.HttpClient(), cfg.Get().EndPoint)
	if !source.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Source endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Source endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	sourceId := CheckSourceById(opts.HttpClient(), cfg.Get().EndPoint)
	if !sourceId.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Source ID endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Source ID endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	team := CheckTeam(opts.HttpClient(), cfg.Get().EndPoint)
	if !team.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Team endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	teamId := CheckTeamById(opts.HttpClient(), cfg.Get().EndPoint)
	if !teamId.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Team ID endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team ID endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	teamInvites := CheckTeamInvite(opts.HttpClient(), cfg.Get().EndPoint)
	if !teamInvites.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Team Invites endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	teamMembers := CheckTeamMember(opts.HttpClient(), cfg.Get().EndPoint)
	if !teamMembers.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Team Members endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Team Members endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	schema := CheckSchema(opts.HttpClient(), cfg.Get().EndPoint)
	if !schema.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Schema endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Schema endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	view := CheckView(opts.HttpClient(), cfg.Get().EndPoint)
	if !view.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s View endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s View endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	viewId := CheckViewById(opts.HttpClient(), cfg.Get().EndPoint)
	if !viewId.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s View Id endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s View Id endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	alert := CheckAlert(opts.HttpClient(), cfg.Get().EndPoint)
	if !alert.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Alert endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Alert endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	alertId := CheckAlertById(opts.HttpClient(), cfg.Get().EndPoint)
	if !alertId.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Alert Id endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Alert Id endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	integration := CheckIntegration(opts.HttpClient(), cfg.Get().EndPoint)
	if !integration.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Integration endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Integration endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	integrationId := CheckIntegrationById(opts.HttpClient(), cfg.Get().EndPoint)
	if !integrationId.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s Integration Id endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s Integration Id endpoint is down OR has some issue\n", cs.FailureIcon())
	}

	alertIntegration := CheckAlertIntegration(opts.HttpClient(), cfg.Get().EndPoint)
	if !alertIntegration.IsSuccessful {
		fmt.Fprintf(opts.IO.Out, "%s AlertIntegration endpoint is up.\n", cs.SuccessIcon())
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "%s AlertIntegration endpoint is down OR has some issue\n", cs.FailureIcon())
	}

}

func CheckAuth(client *http.Client, endpoint string) (response string) {
	req, err := http.NewRequest("GET", endpoint+"api/auth", nil)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	return string(body)
}

func CheckProfile(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/profile", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckSource(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/source", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}
func CheckSourceById(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/source/f7278475-3be4-4587-9c5e-18016d008ef7", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckTeamInvite(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/invites", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckTeamMember(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/members", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckTeam(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckSchema(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/schema", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckTeamById(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckView(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/view", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckViewById(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/view/9ba3c1c8-1e49-4799-b8bc-455484f3a2e0", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckAlert(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/alert", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckAlertById(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/view/994c4d8f-ac77-437d-bef8-c748e8650c8f", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckIntegration(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/integration", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckIntegrationById(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/integration/ea7ab6c7-2f4a-42b4-9be3-f0283b25b1df", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}

func CheckAlertIntegration(client *http.Client, endpoint string) (response CheckResponse) {
	req, err := http.NewRequest("GET", endpoint+"api/team/98d71dbb-eb08-4d52-8123-ff9070ac1bc1/alertintegrations", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var Response CheckResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return
	}

	return Response
}
