package main

import (
	"fmt"
	"log"
	"strings"

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
	field := NewLongAnswerField()
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
	// represents active screen
	currentView string

	questions []Question
	width int
	height int
	index int
	styles *Styles
	done bool

	// Database connection details
	dbType string
	host string
	port int
	user string
	password string
	dbName string
}

func New(questions []Question) *model {
	styles := DefaultStyles()
	return &model{
		currentView: "connect",
		questions: questions,
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
			if m.index == len(m.questions) - 1 {
				m.done = true
			}
			current.answer = current.input.Value()
			log.Printf("question: %s, answer: %s", current.question, current.answer)
			m.Next()
			return m, current.input.Blur
		}		
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.currentView {
	case "connect":
		return renderConnectScreen(m)
	default:
		current := m.questions[m.index]
		if m.done {
			var output string
			for _, q := range m.questions {
				output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
			}
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, lipgloss.JoinHorizontal(lipgloss.Center, lipgloss.NewStyle().Bold(true).Render("Done!"), lipgloss.NewStyle().Bold(false).Render(" Press Ctrl+C to exit.")), output))
		}
		if m.width == 0 {
			return "loading..."
		}
		return lipgloss.Place(
			m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, m.questions[m.index].question, m.styles.InputField.Render(current.input.View())),
		)
	}
}

func renderConnectScreen(m model) string{
	borderWidth := 2
	borderHeight := 2
	// Calculate 1/3 and 2/3 widths
	oneThird := (m.width - borderWidth) / 3
	twoThirds := m.width - oneThird - 2 * borderWidth

	// Calculate padding to center the text vertically
	paddingTop := (m.height - borderHeight - 2) / 2 // -2 because of the top and bottom border
	paddingBottom := m.height - borderHeight - paddingTop - 1 // -1 for the actual text line

	// Style for the left box
	leftBoxStyle := lipgloss.NewStyle().
		Width(oneThird).
		Height(m.height - borderHeight).
		PaddingTop(paddingTop).
		PaddingBottom(paddingBottom).
		Align(lipgloss.Center).
		BorderForeground(lipgloss.Color("#F2CC8F")).
		BorderStyle(lipgloss.ThickBorder())

	// Style for the right box
	rightBoxStyle := lipgloss.NewStyle().
		Width(twoThirds).
		Height(m.height - borderHeight).
		Align(lipgloss.Center).
		BorderForeground(lipgloss.Color("#F2CC8F")).
		BorderStyle(lipgloss.ThickBorder())

	// Rendered boxes
	leftBox := leftBoxStyle.Render("hi there bestie")
	rightBox := rightBoxStyle.Render("we must go BIGGER")

	// Split each box into its lines
	leftLines := strings.Split(leftBox, "\n")
	rightLines := strings.Split(rightBox, "\n")

	// Combine the two boxes line by line
	var combinedLines []string
	for i := 0; i < m.height; i++ {
		combinedLines = append(combinedLines, leftLines[i]+rightLines[i])
	}

	return strings.Join(combinedLines, "\n")
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
		newShortQuestion("what is your name?"),
		newLongQuestion("what is your quest?"),
		newLongQuestion("what is your favorite color?"),
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