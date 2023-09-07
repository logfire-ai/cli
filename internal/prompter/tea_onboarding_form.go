package prompter

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/grpcutil"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

func NewOnboardingForm() {
	p := tea.NewProgram(
		initialModel(),
		// tea.WithMouseCellMotion(),
	)

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
	role
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
	logfireGray = lipgloss.Color("#285cde")

	colorOneStyle = lipgloss.Color("#202f55")

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
	roleList    list.Model
	choice      string
	quitting    bool
	config      config.Config
	log         string
	sourceId    string
	sourceToken string
	content     string
	ready       bool
	viewport    viewport.Model
}

func (m *model) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
}

func initialModel() *model {
	content, err := os.ReadFile("./internal/prompter/artichoke.md")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

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

	roles := []list.Item{
		item("Software Engineer"),
		item("Data Scientist"),
		item("Product Manager"),
	}

	const defaultWidth = 100
	const listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = fmt.Sprintf("%s %s", "Select the source type you want to create:", continueStyle.Render("Choose the environment from which you'll be sending logs"))
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitileStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	var inputs []textinput.Model = make([]textinput.Model, 14)
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

	lr := list.New(roles, itemDelegate{}, 35, 7)
	lr.Title = "What is your Role?"
	lr.SetShowStatusBar(false)
	lr.SetFilteringEnabled(false)
	lr.SetShowHelp(false)
	lr.Styles.Title = listTitileStyle
	lr.Styles.PaginationStyle = paginationStyle

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

	m.viewport.Init()

	return &model{
		inputs:   inputs,
		focused:  0,
		err:      nil,
		stage:    "email",
		spinner:  s,
		list:     l,
		roleList: lr,
		choice:   "",
		config:   cfg,
		log:      "",
		sourceId: "",
		content:  string(content),
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
		m.err = errors.New("we apologize for the inconvenience. There seems to be an error on our end or with our server.\nPlease try again later or contact our support team for assistance")
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
					os.Exit(0)
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
				subStep = "role"
				m.nextInput()
			}
		case "role":
			i, _ := m.roleList.SelectedItem().(item)
			if string(i) != "" {
				err := APICalls.OnboardingFlow(m.config.Get().ProfileID, m.config.Get().Token, m.config.Get().EndPoint,
					m.inputs[firstName].Value(), m.inputs[lastName].Value(), string(i))
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

				err = m.config.UpdateConfig(nil, nil, nil, nil, nil, &team.ID, nil, nil, nil)
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

				// configuration, err := APICalls.GetConfiguration(m.config.Get().Token, m.config.Get().EndPoint, m.config.Get().TeamId, m.sourceId)
				// if err != nil {
				// 	log.Fatal(err)
				// }

				// prettyJSON, err := json.MarshalIndent(configuration, "", "  ")
				// if err != nil {
				// 	log.Fatalf("Error encoding JSON: %s", err)
				// }

				// m.content = string(prettyJSON)

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
		os.Exit(0)
		return m, tea.Quit
	}

	return nil, nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

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

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.focused != 4 && m.focused != 7 && m.focused != 9 && m.focused != 10 && m.focused != 11 && m.focused != 12 && m.focused != 13 {
		for i := range m.inputs {
			m.inputs[i].Blur()
			m.inputs[m.focused].Focus()
			var inputCmd tea.Cmd
			m.inputs[i], inputCmd = m.inputs[i].Update(msg)
			cmds = append(cmds, inputCmd)
		}
	} else if subStep == "role" {
		// Update list
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		var listCmd tea.Cmd
		m.roleList, listCmd = m.roleList.Update(msg)
		cmds = append(cmds, listCmd)
	} else if m.focused == 7 {
		// Update list
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		var listCmd tea.Cmd
		m.list, listCmd = m.list.Update(msg)
		cmds = append(cmds, listCmd)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) renderCurlCommand() string {
	return fmt.Sprintf(
		`
%s
%s %s %s \
%s %s \
%s %s \
%s %s

%s

%s %s %s %s %s
%s
`,
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
}

func renderWelcome() string {
	welcomeText := continueStyle.Render("\nWelcome to logfireAI. We'll guide you through a series of steps to get started with using our platform.\n" +
		"If all goes well, you'd have sent and received your very first log in a couple of minutes.")

	return fmt.Sprintf("%s\n", welcomeText)
}

func renderSection(title string, checked bool) string {
	mark := checkMark
	if !checked {
		mark = unCheckMarked
	}
	return fmt.Sprintf("\n%s%s\n", mark, colorOneBackgroundColorFourForeground.Render(title))
}

func (m model) renderEmail() string {
	return fmt.Sprintf("\n%s\n%s\n", inputStyle.Render("Email"), m.inputs[email].View())
}

func (m model) renderToken() string {
	return fmt.Sprintf("\n%s\n%s\n", inputStyle.Render("Token [paste]"), m.inputs[token].View())
}

func (m model) renderAccountSetup() string {
	return fmt.Sprintf(`
%s   %s
%s   %s

%s
%s
%s
`, inputStyle.Render("First name"),
		inputStyle.Render("Last name"),
		m.inputs[firstName].View(),
		m.inputs[lastName].View(),
		m.roleList.View(),
		inputStyle.Render("Password"),
		m.inputs[password].View())
}

func (m model) renderTeamName() string {
	return fmt.Sprintf("\n%s %s\n%s\n", inputStyle.Render("Team name"), continueStyle.Render("Assign a name for your team; you can interact with one team at a time, each with its own sources and members"),
		m.inputs[teamName].View())
}

func (m model) renderSource() string {
	return fmt.Sprintf(`

%s

%s %s
%s
`, m.list.View(),
		inputStyle.Render("Source name"), continueStyle.Render("Name this source for easy identification among multiple sources in your team"),
		m.inputs[sourceName].View(),
	)
}

func (m model) renderConfig() string {
	// m.viewport.SetContent(m.content)

	if !m.ready {
		return "\n  Initializing..."
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())

	// 	return fmt.Sprintf(`

	// %s

	// %s
	// `, continueStyle.Render("Follow the provided steps to configure your environment for log transmission"),
	//
	//	configuration)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s", m.err)
	}

	waitingForLog := continueStyle.Render(fmt.Sprintf("\nWaiting for logs... %s", m.spinner.View()))

	awesomeLogReceived := colorTwo.Render("\nAwesome! log received\n")

	continueMessage := continueStyle.Render("\nContinue ->")

	switch step {

	case "signup":
		switch subStep {

		case "email":
			return renderWelcome() + renderSection("Signup", false) + m.renderEmail() + continueMessage

		case "token":
			return renderWelcome() + renderSection("Signup", false) + m.renderEmail() + m.renderToken() + continueMessage
		}

	case "account-setup":
		return renderWelcome() + renderSection("Signup", true) + m.renderEmail() + m.renderToken() + renderSection("Account setup", false) + m.renderAccountSetup() + continueMessage

	case "team":
		return renderWelcome() + renderSection("Signup", true) + m.renderEmail() + m.renderToken() + renderSection("Account setup", true) + m.renderAccountSetup() + renderSection("Create a Team", false) + m.renderTeamName() + continueMessage

	case "send-logs":
		switch subStep {
		case "select-source":
			return renderWelcome() + renderSection("Signup", true) + m.renderEmail() + m.renderToken() + renderSection("Account setup", true) + m.renderAccountSetup() + renderSection("Create a Team", true) + m.renderTeamName() + renderSection("Send logs", false) + m.renderSource() + continueMessage

		case "sourceName":
			return renderWelcome() + renderSection("Signup", true) + m.renderEmail() + m.renderToken() + renderSection("Account setup", true) + m.renderAccountSetup() + renderSection("Create a Team", true) + m.renderTeamName() + renderSection("Send logs", false) + m.renderSource() + continueMessage

		case "curl":
			return renderWelcome() + renderSection("Signup", true) + m.renderEmail() + m.renderToken() + renderSection("Account setup", true) + m.renderAccountSetup() + renderSection("Create a Team", true) + m.renderTeamName() + renderSection("Send logs", false) + m.renderSource() + m.renderCurlCommand() + continueMessage

		case "wait":
			return renderWelcome() + renderSection("Signup", true) + m.renderEmail() + m.renderToken() + renderSection("Account setup", true) + m.renderAccountSetup() + renderSection("Create a Team", true) + m.renderTeamName() + renderSection("Send logs", false) + m.renderSource() + m.renderCurlCommand() + waitingForLog

		case "awesome":
			return renderWelcome() + renderSection("Signup", true) + m.renderEmail() + m.renderToken() + renderSection("Account setup", true) + m.renderAccountSetup() + renderSection("Create a Team", true) + m.renderTeamName() + renderSection("Send logs", false) + m.renderSource() + m.renderCurlCommand() + awesomeLogReceived + continueMessage
		}

	case "config-source":
		// return renderWelcome() + renderSection("Signup", true) + m.renderEmail() + m.renderToken() + renderSection("Account setup", true) + m.renderAccountSetup() + renderSection("Create a Team", true) + m.renderTeamName() + renderSection("Send logs", true) + m.renderSource() + m.renderCurlCommand() + awesomeLogReceived + renderSection("Config source", false) +
		return m.renderConfig()
		// + continueMessage

	case "complete":

		return "Completed!"
	}

	return ""
}

func (m model) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
// func (m *model) prevInput() {
// 	m.focused--
// 	// Wrap around
// 	if m.focused < 0 {
// 		m.focused = len(m.inputs) - 1
// 	}
// }
