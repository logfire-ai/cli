package prompter

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	//s.BorderColor = lipgloss.Color("36")
	//s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type Main struct {
	styles    *Styles
	index     int
	Questions []Question
	width     int
	height    int
	done      bool
}

type Question struct {
	question string
	Answer   string
	input    Input
}

type exitMsg struct {
	msg string
}

func NewQuestion(q string) Question {
	return Question{question: q}
}

func NewShortQuestion(q, placeholder string) Question {
	question := NewQuestion(q)
	model := NewShortAnswerField(placeholder)
	question.input = model
	return question
}

func NewLongQuestion(q string) Question {
	question := NewQuestion(q)
	model := NewLongAnswerField(q)
	question.input = model
	return question
}

func NewSelectableQuestion(q string) Question {
	return Question{}
}

func NewTea(Questions []Question) *Main {
	styles := DefaultStyles()
	return &Main{styles: styles, Questions: Questions}
}

func (m Main) Init() tea.Cmd {
	return m.Questions[m.index].input.Blink
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.Questions[m.index]
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.index == len(m.Questions)-1 {
				m.done = true
				current.Answer = current.input.Value()
				return m, tea.Quit
			}
			current.Answer = current.input.Value()
			m.Next()
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m Main) View() string {
	current := m.Questions[m.index]
	if m.done {
		var output string
		for _, q := range m.Questions {
			output += fmt.Sprintf("%s: %s\n", q.question, q.Answer)
		}
		return output
	}
	if m.width == 0 {
		return "loading..."
	}
	// stack some left-aligned strings together in the center of the window
	//return lipgloss.Place(
	//	m.width,
	//	m.height,
	//	lipgloss.Center,
	//	lipgloss.Center,
	//	lipgloss.JoinVertical(
	//		lipgloss.Left,
	//		current.question,
	//		m.styles.InputField.Render(current.input.View()),
	//	),
	//)

	return fmt.Sprintf(
		"%s\n%s\n%s\n",
		current.question,
		m.styles.InputField.Render(current.input.View()),
		"(esc to quit)",
	) + "\n"
}

func (m *Main) Next() {
	if m.index < len(m.Questions)-1 {
		m.index++
	}
}

func Initializer(main *Main) (*Main, error) {
	// init styles; optional, just showing as a way to organize styles
	// start bubble tea and init first model

	var err error

	p := tea.NewProgram(*main)
	if _, err = p.Run(); err != nil {
		log.Fatal(err)
	}

	return main, err
}
