package prompter

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"
)

func NewOnboardingForm() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

const (
	email = iota
	token
	firstName
	lastName
	password
	teamName
	selectSource
	sourceName
	curl
	wait
	awesome
	sourceConfig
	Complete
)

const (
	darkGray    = lipgloss.Color("#767676")
	logfireGray = lipgloss.Color("#607d8b")

	colorOneStyle = lipgloss.Color("#27374D")

	colorTwoStyle = lipgloss.Color("#526D82")

	colorThreeStyle = lipgloss.Color("#9DB2BF")

	colorFourStyle             = lipgloss.Color("#DDE6ED")
	stepText                   = lipgloss.Color("#DDE6ED")
	textHiglight               = lipgloss.Color("#0089d2")
	commandBackgroundHighlight = lipgloss.Color("#5f5f5f")
)

var (
	colorOneBackgroundColorFourForeground = lipgloss.NewStyle().Background(colorOneStyle).Foreground(colorFourStyle)
	colorOne                              = lipgloss.NewStyle().Foreground(colorOneStyle)
	colorTwo                              = lipgloss.NewStyle().Foreground(colorTwoStyle)
	colorThree                            = lipgloss.NewStyle().Foreground(colorThreeStyle)
	colorFour                             = lipgloss.NewStyle().Foreground(colorFourStyle)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(0)
	listTitileStyle   = lipgloss.NewStyle().Foreground(logfireGray).PaddingLeft(0).MarginLeft(0)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(colorTwoStyle)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(0)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(0).PaddingBottom(0)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	inputStyle        = lipgloss.NewStyle().Foreground(logfireGray)
	continueStyle     = lipgloss.NewStyle().Foreground(darkGray)
	codeStyle         = lipgloss.NewStyle().Foreground(colorFourStyle)
	spinnerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Background(colorOneStyle).Foreground(colorFourStyle).Faint(true).Render("[✓] ")
	unCheckMarked       = lipgloss.NewStyle().Background(colorOneStyle).Foreground(colorFourStyle).Faint(true).Render("[ ] ")

	highlightedText                 = lipgloss.NewStyle().Foreground(textHiglight)
	commandBackgroundHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f2f2f2")).Background(commandBackgroundHighlight)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	inputs      []textinput.Model
	spinner     spinner.Model
	focused     int
	err         error
	stage       string
	endpoint    string
	list        list.Model
	choice      string
	quitting    bool
	config      config.Config
	log         string
	sourceId    string
	sourceToken string
}

func (m *model) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
}

func initialModel() *model {
	cfg, _ := config.NewConfig()

	m := &model{}

	m.resetSpinner()

	items := []list.Item{
		item("Kubernetes"),
		item("AWS"),
		item("JavaScript"),
		item("Docker"),
		item("Nginx"),
		item("Dokku"),
		item("Fly.io"),
		item("Heroku"),
		item("Ubuntu"),
		item("Vercel"),
		item(".Net"),
		item("Apache2"),
		item("Cloudflare"),
		item("Java"),
		item("Python"),
		item("PHP"),
		item("PostgreSQL"),
		item("Redis"),
		item("Ruby"),
		item("Mongodb"),
		item("MySQL"),
		item("HTTP"),
		item("Vector"),
		item("fluentbit"),
		item("Fluentd"),
		item("Logstash"),
		item("Rsyslog"),
		item("Render"),
		item("syslog-ng"),
	}

	const defaultWidth = 100
	const listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select the source type you want to create:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitileStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	var inputs []textinput.Model = make([]textinput.Model, 13)
	inputs[email] = textinput.New()
	inputs[email].Placeholder = "Enter your email address here"
	inputs[email].Focus()
	inputs[email].Width = 30
	inputs[email].Prompt = ""

	inputs[token] = textinput.New()
	inputs[token].Placeholder = "Token"
	inputs[token].Focus()
	inputs[token].Width = 100
	inputs[token].Prompt = ""

	inputs[firstName] = textinput.New()
	inputs[firstName].Placeholder = "First name"
	inputs[firstName].Width = 9
	inputs[firstName].Prompt = ""

	inputs[lastName] = textinput.New()
	inputs[lastName].Placeholder = "Last name"
	inputs[lastName].Width = 10
	inputs[lastName].Prompt = ""

	inputs[password] = textinput.New()
	inputs[password].Placeholder = "Password"
	inputs[password].Width = 10
	inputs[password].Prompt = ""
	inputs[password].EchoMode = textinput.EchoPassword
	inputs[password].EchoCharacter = '•'

	inputs[teamName] = textinput.New()
	inputs[teamName].Placeholder = "Team Name"
	inputs[teamName].Width = 8
	inputs[teamName].Prompt = ""

	inputs[sourceName] = textinput.New()
	inputs[sourceName].Placeholder = "Source Name"
	inputs[sourceName].Width = 20
	inputs[sourceName].Prompt = ""

	return &model{
		inputs:   inputs,
		focused:  0,
		err:      nil,
		stage:    "email",
		spinner:  s,
		list:     l,
		choice:   "",
		config:   cfg,
		log:      "",
		sourceId: "",
	}
}

func (m model) Init() tea.Cmd {
	var cmds = []tea.Cmd{textinput.Blink, m.spinner.Tick}

	return tea.Batch(cmds...)
}

var sourceCreated bool

var stop = make(chan error)

var step = "signup"

var subStep = "email"

func waitForLog(m *model) {
	go grpcutil.GetLog(m.config.Get().Token, m.config.Get().EndPoint, m.config.Get().TeamId, m.sourceId, stop)
	err := <-stop
	if err != nil {
		m.err = errors.New("We apologize for the inconvenience. There seems to be an error on our end or with our server.\nPlease try again later or contact our support team for assistance.")
		m.nextInput()
	}
	subStep = "awesome"
	m.nextInput()
}

func (m *model) handleKeyPres() (tea.Model, tea.Cmd) {
	switch step {
	case "signup":
		switch subStep {
		case "email":
			if m.inputs[email].Value() != "" {
				msg, err := APICalls.SignupFlow(m.inputs[email].Value(), m.config.Get().EndPoint)
				if err != nil {

					m.err = err
					return m, nil
				} else if msg == "already registered user. Sent link to login" {
					m.err = errors.New("you are already a user, please use logfire commands")
					return m, nil
				}
				subStep = "token"
				m.nextInput()
			}
		case "token":
			if m.inputs[token].Value() != "" {
				err := APICalls.TokenSignIn(m.config, m.inputs[token].Value(), m.config.Get().EndPoint)
				if err != nil {
					m.err = err
					return m, nil
				}
				step = "account-setup"
				subStep = "firstName"
				m.nextInput()
			}
		}
	case "account-setup":
		switch subStep {
		case "firstName":
			if m.inputs[firstName].Value() != "" {
				subStep = "lastName"
				m.nextInput()
			}
		case "lastName":
			if m.inputs[lastName].Value() != "" {
				err := APICalls.OnboardingFlow(m.config.Get().ProfileID, m.config.Get().Token, m.config.Get().EndPoint,
					m.inputs[firstName].Value(), m.inputs[lastName].Value())
				if err != nil {
					m.err = err
					return m, nil
				}
				subStep = "password"
				m.nextInput()
			}
		case "password":
			if m.inputs[password].Value() != "" {
				err := APICalls.SetPassword(m.config.Get().Token, m.config.Get().EndPoint, m.config.Get().ProfileID,
					m.inputs[password].Value())
				if err != nil {
					m.err = err
					return m, nil
				}
				step = "team"
				subStep = "teamName"
				m.nextInput()
			}
		}
	case "team":
		switch subStep {
		case "teamName":
			if m.inputs[teamName].Value() != "" {
				team, err := APICalls.CreateTeam(m.config.Get().Token, m.config.Get().EndPoint,
					m.inputs[teamName].Value())
				if err != nil {
					m.err = err
					return m, nil
				}

				err = m.config.UpdateConfig(nil, nil, nil, nil, &team.ID, nil, nil, nil)
				if err != nil {
					m.err = err
					return m, nil
				}

				step = "send-logs"
				subStep = "select-source"
				m.nextInput()
			}
		}
	case "send-logs":
		switch subStep {
		case "select-source":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			subStep = "sourceName"
			m.nextInput()
		case "sourceName":
			if m.inputs[sourceName].Value() != "" && !sourceCreated {
				source, err := APICalls.CreateSource(m.config.Get().Token, m.config.Get().EndPoint, m.config.Get().TeamId, m.inputs[sourceName].Value(), m.choice)
				if err != nil {
					m.err = err
					return m, nil
				}

				m.sourceId = source.ID
				m.sourceToken = source.SourceToken
				subStep = "curl"
				m.nextInput()
			}
		case "curl":
			subStep = "wait"
			m.nextInput()

			time.AfterFunc(time.Second*3, func() {
				waitForLog(m)
			})
		case "wait":

		case "awesome":
			step = "config-source"
			subStep = "source-config"
			m.nextInput()
		}
	case "config-source":
		switch subStep {
		case "source-config":
			step = "complete"
			m.nextInput()
		}
	case "complete":
		return m, tea.Quit
	}

	return nil, nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.handleKeyPres()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		//case tea.KeyShiftTab, tea.KeyCtrlP:
		//	m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.handleKeyPres()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.focused != 6 && m.focused != 8 && m.focused != 9 && m.focused != 10 && m.focused != 11 && m.focused != 12 {
		for i := range m.inputs {
			m.inputs[i].Blur()
			m.inputs[m.focused].Focus()
		}
	} else {
		// Update list
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		var listCmd tea.Cmd
		m.list, listCmd = m.list.Update(msg)
		cmds = append(cmds, listCmd)
	}

	// Update inputs
	for i := range m.inputs {
		var inputCmd tea.Cmd
		m.inputs[i], inputCmd = m.inputs[i].Update(msg)
		cmds = append(cmds, inputCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s", m.err)
	}

	curlCommand := fmt.Sprintf(

		`%s
%s %s %s \
%s %s \
%s %s \
%s %s

%s

%s %s %s %s %s
%s`,
		continueStyle.Render("******************************************************************************************"),
		colorThree.Render(`curl`),
		colorThree.Render(`--location`),
		colorThree.Render(`'https://in.logfire.ai'`),
		colorThree.Render(`--header`),
		colorThree.Render(`'Content-Type: application/json'`),
		colorThree.Render(`--header`),
		colorThree.Render(`'Authorization: Bearer `+m.sourceToken+`'`),
		colorThree.Render(`--data`),
		colorThree.Render(`'[{"dt":"2023-06-15T6:00:39.351Z","message":"Hello from Logfire fluentd!"}]'`),
		colorTwo.Render("\nCopy the command and paste it in another terminal to test the source"),
		`Press`, commandBackgroundHighlightStyle.Render(`Enter`),
		`or`,
		commandBackgroundHighlightStyle.Render(`Tab`),
		`after you have copied`,
		continueStyle.Render("******************************************************************************************"),
	)

	waitingForLog := continueStyle.Render("Waiting for logs...")

	awesomeLogReceived := colorTwo.Render("Awesome! log received")

	welcomeTest := continueStyle.Render("Welcome to logfireAI. We'll guide you through a series of steps to get started with using our platform.\n" +
		"If all goes well, you'll have sent and received your very first log in a couple of minutes.")

	switch step {
	case "signup":
		switch subStep {
		case "email":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s
	
%s
	`,
				welcomeTest,
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				continueStyle.Render("Continue ->"),
			) + "\n"

		case "token":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s
	`,
				welcomeTest,
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				continueStyle.Render("Continue ->"),
			) + "\n"
		}
	case "account-setup":
		switch subStep {
		case "firstName":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s
	`,
				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				continueStyle.Render("Continue ->"),
			) + "\n"
		case "lastName":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s
	`,
				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				continueStyle.Render("Continue ->"),
			) + "\n"
		case "password":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s
	`,
				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				continueStyle.Render("Continue ->"),
			) + "\n"
		}
	case "team":
		switch subStep {
		case "teamName":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s%s

%s
%s

%s
	`,

				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Team name"),
				inputStyle.Render("Team name"),
				m.inputs[teamName].View(),
				continueStyle.Render("Continue ->"),
			) + "\n"
		}
	case "send-logs":
		switch subStep {
		case "select-source":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s%s

%s
%s

%s%s

%s

%s
%s

%s
	`,

				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Team name"),
				inputStyle.Render("Team name"),
				m.inputs[teamName].View(),
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Send logs"),
				m.list.View(),
				inputStyle.Render("Source name"),
				m.inputs[sourceName].View(),
				continueStyle.Render("Continue ->"),
			) + "\n"
		case "sourceName":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s%s

%s
%s

%s%s

%s

%s
%s

%s
	`,

				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Team name"),
				inputStyle.Render("Team name"),
				m.inputs[teamName].View(),
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Send logs"),
				m.list.View(),
				inputStyle.Render("Source name"),
				m.inputs[sourceName].View(),
				continueStyle.Render("Continue ->"),
			) + "\n"
		case "curl":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s%s

%s
%s

%s%s

%s

%s
%s

%s

%s
	`,

				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Team name"),
				inputStyle.Render("Team name"),
				m.inputs[teamName].View(),
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Send logs"),
				m.list.View(),
				inputStyle.Render("Source name"),
				m.inputs[sourceName].View(),

				curlCommand,

				continueStyle.Render("Continue ->"),
			) + "\n"
		case "wait":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s%s

%s
%s

%s%s

%s

%s
%s

%s

%s %s
	`,

				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Team name"),
				inputStyle.Render("Team name"),
				m.inputs[teamName].View(),
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Send logs"),
				m.list.View(),
				inputStyle.Render("Source name"),
				m.inputs[sourceName].View(),

				curlCommand,
				waitingForLog,
				m.spinner.View(),
			) + "\n"
		case "awesome":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s%s

%s
%s

%s%s

%s

%s
%s

%s

%s%s

%s
	`,

				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Team name"),
				inputStyle.Render("Team name"),
				m.inputs[teamName].View(),
				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Send logs"),
				m.list.View(),
				inputStyle.Render("Source name"),
				m.inputs[sourceName].View(),

				curlCommand,
				awesomeLogReceived,
				" ",
				continueStyle.Render("Continue ->"),
			) + "\n"
		}
	case "config-source":
		switch subStep {
		case "source-config":
			return fmt.Sprintf(
				`%s

%s%s
	
%s 
%s

%s 
%s
	
%s%s

%s  %s
%s  %s

%s 
%s

%s%s

%s
%s

%s%s

%s

%s
%s

%s

%s%s

%s%s

%s
			
%s
				`,
				welcomeTest,
				checkMark, colorOneBackgroundColorFourForeground.Render("Signup"),
				inputStyle.Render("Email"),
				m.inputs[email].View(),
				inputStyle.Render("Token [paste]"),
				m.inputs[token].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Account setup"),
				inputStyle.Render("First name"),
				inputStyle.Render("Last name"),
				m.inputs[firstName].View(),
				m.inputs[lastName].View(),
				inputStyle.Render("Password"),
				m.inputs[password].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Team name"),
				inputStyle.Render("Team name"),
				m.inputs[teamName].View(),
				checkMark, colorOneBackgroundColorFourForeground.Render("Send logs"),
				m.list.View(),
				inputStyle.Render("Source name"),
				m.inputs[sourceName].View(),

				curlCommand,
				awesomeLogReceived,
				" ",

				unCheckMarked, colorOneBackgroundColorFourForeground.Render("Config source"),

				`You can setup your source using this config`,

				continueStyle.Render("Continue ->"),
			) + "\n"
		}
	case "complete":
		return "Completed!"
	}

	return ""
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
