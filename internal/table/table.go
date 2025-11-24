package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/presselam/yadc/internal/bubble"
	"log"
)

type contextState uint

const (
	None             = iota
	ImageContext     = iota
	ContainerContext = iota
	VolumeContext    = iota
	InspectContext   = iota
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	table   bubble.Model
	width   int
	context contextState
}

type action func(*Model, string)

type KeyMapping struct {
	key key.Binding
	cmd action
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
		m.actionHandler(msg)
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

	var err error
	m.context = context
	s := m.table.Styles()
	switch m.context {
	case ContainerContext:
		err = m.PopulateContainers()
		s.Cell = ContainerFormatter
		m.table.SetStyles(s)
	case ImageContext:
		err = m.PopulateImages()
		s.Cell = nil
		m.table.SetStyles(s)
	}

	return err
}

func (m *Model) actionHandler(msg tea.KeyMsg) {
	var mappings []KeyMapping
	switch m.context {
	case ContainerContext:
		mappings = m.containerKeyMapping()
	case ImageContext:
		m.imageHandler(msg)
	}

	row := m.table.SelectedRow()
	for _, command := range mappings {
		if key.Matches(msg, command.key) {
			command.cmd(m, row[0])
		}
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

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	m := Model{
		table: t,
	}
	m.SetContext(ContainerContext)

	return m
}
