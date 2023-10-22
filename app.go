package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const SelectedColor = lipgloss.Color("#01BE85")

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

	selectedTabIdx int
	totalTabs int
	activeTab string

	selectedRowIdx int
	selectedRowData []string

	data [][] string

	table *table.Table

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
	data := [][]string{
		{"1", "Bulbasaur", "Grass", "Poison", "フシギダネ", "Bulbasaur"},
		{"2", "Ivysaur", "Grass", "Poison", "フシギソウ", "Ivysaur"},
		{"3", "Venusaur", "Grass", "Poison", "フシギバナ", "Venusaur"},
		{"4", "Charmander", "Fire", "", "ヒトカゲ", "Hitokage"},
		{"5", "Charmeleon", "Fire", "", "リザード", "Lizardo"},
		{"6", "Charizard", "Fire", "Flying", "リザードン", "Lizardon"},
		{"7", "Squirtle", "Water", "", "ゼニガメ", "Zenigame"},
		{"8", "Wartortle", "Water", "", "カメール", "Kameil"},
		{"9", "Blastoise", "Water", "", "カメックス", "Kamex"},
		{"10", "Caterpie", "Bug", "", "キャタピー", "Caterpie"},
		{"11", "Metapod", "Bug", "", "トランセル", "Trancell"},
		{"12", "Butterfree", "Bug", "Flying", "バタフリー", "Butterfree"},
		{"13", "Weedle", "Bug", "Poison", "ビードル", "Beedle"},
		{"14", "Kakuna", "Bug", "Poison", "コクーン", "Cocoon"},
		{"15", "Beedrill", "Bug", "Poison", "スピアー", "Spear"},
		{"16", "Pidgey", "Normal", "Flying", "ポッポ", "Poppo"},
		{"17", "Pidgeotto", "Normal", "Flying", "ピジョン", "Pigeon"},
		{"18", "Pidgeot", "Normal", "Flying", "ピジョット", "Pigeot"},
		{"19", "Rattata", "Normal", "", "コラッタ", "Koratta"},
		{"20", "Raticate", "Normal", "", "ラッタ", "Ratta"},
		{"21", "Spearow", "Normal", "Flying", "オニスズメ", "Onisuzume"},
		{"22", "Fearow", "Normal", "Flying", "オニドリル", "Onidrill"},
		{"23", "Ekans", "Poison", "", "アーボ", "Arbo"},
		{"24", "Arbok", "Poison", "", "アーボック", "Arbok"},
		{"25", "Pikachu", "Electric", "", "ピカチュウ", "Pikachu"},
		{"26", "Raichu", "Electric", "", "ライチュウ", "Raichu"},
		{"27", "Sandshrew", "Ground", "", "サンド", "Sand"},
		{"28", "Sandslash", "Ground", "", "サンドパン", "Sandpan"},
	}
	return &model{
		data: data,
		selectedTabIdx: -1,
		totalTabs: 2,
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
		case "tab":
			m.Next()
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch msg.String() {
				case "esc":
					m.selectedRowIdx = 0
					m.activeTab = "none"
				case "enter":
					log.Print(m.data[m.selectedRowIdx])
					log.Printf("TESTING")
					m.selectedRowData = m.data[m.selectedRowIdx]
					m.activeTab = "table"
					updateTable(&m)
				case "up":
					if m.selectedRowIdx > 0 {
						m.selectedRowIdx--
					}
				case "down":
					if m.selectedRowIdx < len(m.selectedRowData) - 1 {
						m.selectedRowIdx++
					}
				}
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
	paddingTop := (m.height - borderHeight - 2) / 2 - 10 // -2 because of the top and bottom border
	paddingBottom := m.height - borderHeight - paddingTop - 1 - 10 // -1 for the actual text line

	borderColor := lipgloss.Color("#F2CC8F")
	if m.selectedTabIdx == 0 {
		borderColor = SelectedColor
	}

	// Style for the left box
	leftBoxStyle := lipgloss.NewStyle().
		Width(oneThird).
		Height(m.height - borderHeight - 100).
		PaddingTop(paddingTop).
		PaddingBottom(paddingBottom).
		Align(lipgloss.Center).
		BorderForeground(borderColor).
		BorderStyle(lipgloss.ThickBorder())

	// Style for the right box
	// rightBoxStyle := lipgloss.NewStyle().
	// 	Width(twoThirds).
	// 	Height(m.height - borderHeight).
	// 	Align(lipgloss.Center).
	// 	BorderForeground(lipgloss.Color("#F2CC8F")).
	// 	BorderStyle(lipgloss.ThickBorder())

	// Rendered boxes
	leftBox := leftBoxStyle.Render("hi there bestie")
	// rightBox := rightBoxStyle.Render("we must go BIGGER")

	// Split each box into its lines
	// leftLines := strings.Split(leftBox, "\n")
	// rightLines := strings.Split(rightBox, "\n")

	// Combine the two boxes line by line
	// var combinedLines []string
	// for i := 0; i < m.height; i++ {
	// 	combinedLines = append(combinedLines, leftLines[i]+rightLines[i])
	// }

	var table = genTable(&m, m.height, twoThirds)
	m.table = table

	return lipgloss.JoinVertical(lipgloss.Center, lipgloss.JoinHorizontal(lipgloss.Center, leftBox, table.Render()))
}

func genTable(m *model, height int, width int) *table.Table {
	re := lipgloss.NewRenderer(os.Stdout)
	selectedRowIdx := 1
	selectedRowIdx = m.selectedRowIdx
	log.Printf("selectedRowIdx: %d", selectedRowIdx)
	// selectedStyle := baseStyle.Copy().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
	borderColor := lipgloss.Color("#F2CC8F")
	if m.selectedTabIdx == 1 {
		borderColor = SelectedColor
	}
	if m.activeTab == "table" {
		borderColor = SelectedColor
		selectedRowIdx = 0
	}
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Copy().Foreground(lipgloss.Color("252")).Bold(true)
	typeColors := map[string]lipgloss.Color{
		"Bug":      lipgloss.Color("#D7FF87"),
		"Electric": lipgloss.Color("#FDFF90"),
		"Fire":     lipgloss.Color("#FF7698"),
		"Flying":   lipgloss.Color("#FF87D7"),
		"Grass":    lipgloss.Color("#75FBAB"),
		"Ground":   lipgloss.Color("#FF875F"),
		"Normal":   lipgloss.Color("#929292"),
		"Poison":   lipgloss.Color("#7D5AFC"),
		"Water":    lipgloss.Color("#00E2C7"),
	}
	dimTypeColors := map[string]lipgloss.Color{
		"Bug":      lipgloss.Color("#97AD64"),
		"Electric": lipgloss.Color("#FCFF5F"),
		"Fire":     lipgloss.Color("#BA5F75"),
		"Flying":   lipgloss.Color("#C97AB2"),
		"Grass":    lipgloss.Color("#59B980"),
		"Ground":   lipgloss.Color("#C77252"),
		"Normal":   lipgloss.Color("#727272"),
		"Poison":   lipgloss.Color("#634BD0"),
		"Water":    lipgloss.Color("#439F8E"),
	}

	headers := []string{"#", "Name", "Type 1", "Type 2", "Japanese", "Official Rom."}

	CapitalizeHeaders := func(data []string) []string {
		for i := range data {
			data[i] = strings.ToUpper(data[i])
		}
		return data
	}
	

	t := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(re.NewStyle().Foreground(borderColor)).
		Headers(CapitalizeHeaders(headers)...).
		Width(width).
		Height(height).
		Rows(m.data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}

			even := row%2 == 0

			switch col {
			case 2, 3: // Type 1 + 2
				c := typeColors
				if even {
					c = dimTypeColors
				}

				color := c[fmt.Sprint(m.data[row-1][col])]
				return baseStyle.Copy().Foreground(color)
			}

			if even {
				return baseStyle.Copy().Foreground(lipgloss.Color("245"))
			}
			return baseStyle.Copy().Foreground(lipgloss.Color("252"))
		})

	
	return t
}

func updateTable(m *model) {
	// NEED HELP HERE???
	log.Print("updating table")
	m.table.StyleFunc(func(row, col int) lipgloss.Style {
		if row == m.selectedRowIdx {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	})
		
}

func (m *model) Next() {
	if m.selectedTabIdx < m.totalTabs - 1 {
		m.selectedTabIdx++
	} else {
		m.selectedTabIdx = 0
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