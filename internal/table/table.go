package table

import (
	//	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/presselam/yadc/internal/bubble"
	"log"
)

type contextState uint

const (
	None       = iota
	Images     = iota
	Containers = iota
	Volumes    = iota
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	table   bubble.Model
	width   int
	context contextState
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "enter":
			return m, tea.Batch(
				tea.Printf("selected: [%s]", m.table.SelectedRow()[1]),
			)
		default:
			m.actionHandler(msg.String())
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return baseStyle.Render(m.table.View())
}

func (m *Model) resize(width int, height int) {
	log.Println("table.resize(", width, ",", height, ")")
	m.width = width - 4
	m.table.SetWidth(m.width)
	m.table.SetHeight(height - 12)
}

func (m *Model) SetContext(context contextState) error {
	log.Printf("context:[%v]\n", context)

	m.context = context
	m.populateTable()

	return nil
}

func (m *Model) populateTable() {
	switch m.context {
	case Containers:
		m.PopulateContainers()
	case Images:
		m.PopulateImages()
	}
}

func (m *Model) actionHandler(command string) {
	log.Printf("table.action.%s", command)
	switch m.context {
	case Containers:
		m.containerHandler(command)
	case Images:
		m.imageHandler(command)
	}
}

func New() Model {

	t := bubble.New(
		bubble.WithFocused(true),
	)

	s := bubble.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(false).
		Bold(true)

	s.Cell = s.Cell.
		BorderForeground(lipgloss.Color("40")).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	m := Model{
		table:   t,
		context: Containers,
	}
	m.populateTable()

	return m
}
