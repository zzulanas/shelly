package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField lipgloss.Style
}

type Question struct {
	question string
	answer string
	input Input
}

func NewQuestion(question string) Question {
	return Question{question: question}
}

func newShortQuestion(question string) Question {
	q := NewQuestion(question)
	field := NewShortAnswerField()
	q.input = field
	return q
}

func newLongQuestion(question string) Question {
	q := NewQuestion(question)
	field := NewShortAnswerField()
	q.input = field
	return q
}


func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#F2CC8F")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.RoundedBorder()).Padding(1).Width(80)

	return s
}


type model struct {
	questions []Question
	width int
	height int
	answerField textinput.Model
	index int
	styles *Styles
}

func New(questions []Question) *model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = "Your answer here"
	answerField.Focus()
	return &model{
		questions: questions,
		answerField: answerField,
		styles: styles,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			current.answer = m.answerField.Value()
			m.answerField.SetValue("")
			log.Printf("question: %s, answer: %s", current.question, current.answer)
			m.Next()
			return m, nil
		}		
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, m.questions[m.index].question, m.styles.InputField.Render(m.answerField.View())),
	)
	
}

func (m *model) Next() {
	if m.index < len(m.questions) - 1 {
		m.index++
	} else {
		m.index = 0
	}
}

func main() {
	questions := []Question{
		NewQuestion("what is your name?"),
		NewQuestion("what is your quest?"),
		NewQuestion("what is your favorite color?"),
	}
	m := New(questions)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	p:=tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}