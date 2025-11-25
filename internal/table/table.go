package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/presselam/yadc/internal/bubble"
	"github.com/presselam/yadc/internal/timers"
	"log"
	"sort"
	"time"
)

type ContextState uint

const (
	None             = iota
	ImageContext     = iota
	ContainerContext = iota
	VolumeContext    = iota
	InspectContext   = iota
	LogsContext      = iota
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	id       int
	table    bubble.Model
	width    int
	context  ContextState
	selected string
	sorted   int
}

type action func(*Model, string)

type KeyMapping struct {
	key key.Binding
	cmd action
}

var sortKeys = []key.Binding{
	key.NewBinding(key.WithKeys("1")),
	key.NewBinding(key.WithKeys("2")),
	key.NewBinding(key.WithKeys("3")),
	key.NewBinding(key.WithKeys("4")),
	key.NewBinding(key.WithKeys("5")),
	key.NewBinding(key.WithKeys("1")),
	key.NewBinding(key.WithKeys("1")),
	key.NewBinding(key.WithKeys("1")),
}

func (m Model) tick() tea.Cmd {
	var delay time.Duration

	switch m.context {
	case LogsContext:
		delay = 2 * time.Second
	case InspectContext:
		return nil
	default:
		delay = 5 * time.Second

	}

	log.Println("table.tick:[", delay, "]")
	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return timers.TimerMsg{ID: m.id, Tag: t, Timeout: false}
	})
}

func (m Model) Init() tea.Cmd {
	log.Println("table.init")
	return m.tick()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	log.Printf("table.update: [%v]", msg)

	var cmd tea.Cmd
	var batch []tea.Cmd

	switch msg := msg.(type) {
	case timers.TimerMsg:
		if msg.ID == m.id {
			log.Printf("tick: [%v]", msg)
			// a == a  so that it repopulates the data
			// fix it
			m.SetContext(m.context)
			batch = append(batch, m.tick())
		}
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, m.tick()
		}
		if m.actionHandler(msg) {
			return m, nil
		}
	}

	m.table, cmd = m.table.Update(msg)
	batch = append(batch, cmd)
	return m, tea.Batch(batch...)
}

func (m Model) View() string {
	return baseStyle.Render(m.table.View())
}

func (m *Model) resize(width int, height int) {
	log.Println("table.resize(", width, ",", height, ")")
	m.width = width - 3
	m.table.SetWidth(m.width)
	m.table.SetHeight(height - 9)
}

func (m Model) Context() ContextState {
	return m.context
}

func (m *Model) SetContext(context ContextState) error {
	log.Printf("context:[%v]\n", context)

	var err error
	m.context = context
	s := m.table.Styles()
	switch m.context {
	case ContainerContext:
		err = m.PopulateContainers()
		s.Cell = ContainerFormatter
	case ImageContext:
		err = m.PopulateImages()
		s.Cell = ImageFormatter
	case LogsContext:
		err = m.FetchLogs()
		s.Cell = nil
	case InspectContext:
		s.Cell = nil
	}
	m.table.SetStyles(s)

	return err
}

func (m *Model) actionHandler(msg tea.KeyMsg) bool {
	// check sortkeys
	for i, sortKey := range sortKeys {
		if key.Matches(msg, sortKey) {
			m.sorted = i
			m.sortRows()
			return true
		}
	}

	// check context actions
	var mappings []KeyMapping
	switch m.context {
	case ContainerContext:
		mappings = m.containerKeyMapping()
	case ImageContext:
		mappings = m.imageKeyMapping()
	}

	row := m.table.SelectedRow()
	for _, command := range mappings {
		if key.Matches(msg, command.key) {
			command.cmd(m, row[0])
			return true
		}
	}

	return false
}

func (m *Model) sortRows() {

	rows := m.table.Rows()
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][m.sorted] < rows[j][m.sorted]
	})
	m.table.SetRows(rows)
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
		Foreground(lipgloss.Color("254")).
		Background(lipgloss.Color("12")).
		Bold(false)

	t.SetStyles(s)

	m := Model{
		id:     timers.NextID(),
		table:  t,
		sorted: 1,
	}

	return m
}
