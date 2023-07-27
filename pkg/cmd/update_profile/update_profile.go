package update_profile

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/internal/prompter"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

type UpdateProfileOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	HttpClient func() *http.Client
	Config     func() (config.Config, error)

	Interactive bool

	FirstName string
	LastName  string
}

func UpdateProfileCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UpdateProfileOptions{
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "update-profile",
		Args:  cobra.ExactArgs(0),
		Short: "update your profile",
		Long: heredoc.Docf(`
			update your profile
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ logfire update-profile

			# update profile by reading the details from the prompt
			$ logfire update-profile --first-name <first-name> --last-name <last-name>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}

			UpdateProfileRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.FirstName, "firstname", "", "First name of the user.")
	cmd.Flags().StringVar(&opts.LastName, "lastname", "", "Last name of the user.")
	return cmd
}

func UpdateProfileRun(opts *UpdateProfileOptions) {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read config\n", cs.FailureIcon())
	}

	if opts.Interactive && opts.FirstName == "" && opts.LastName == "" {
		updateFirstName, err := opts.Prompter.Confirm("Do you want to update your First name?", false)

		if updateFirstName {
			opts.FirstName, err = opts.Prompter.Input("Enter First name:", "")
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Name\n", cs.FailureIcon())
				return
			}
		}

		updateLastName, err := opts.Prompter.Confirm("Do you want to update your Last name?", false)

		if updateLastName {
			opts.LastName, err = opts.Prompter.Input("Enter Last name:", "")
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "%s Failed to read Name\n", cs.FailureIcon())
				return
			}
		}

	} else {
		if opts.FirstName == "" {
			fmt.Fprint(opts.IO.ErrOut, "First name is required.")
			os.Exit(0)
		}
	}

	err = APICalls.UpdateProfile(opts.HttpClient(), cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().ProfileID, opts.FirstName, opts.LastName)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Profile updated successfully!\n", cs.SuccessIcon())
	}
}
