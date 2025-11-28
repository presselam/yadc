package monitor

import (
	"errors"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/presselam/yadc/internal/banner"
	"github.com/presselam/yadc/internal/dialog"
	"github.com/presselam/yadc/internal/logger"
	"github.com/presselam/yadc/internal/table"
	"github.com/presselam/yadc/internal/timers"
	"log"
	//	"strconv"
	"strings"
	"time"
)

type sessionState uint

const (
	defaultTime                = time.Minute
	tableFocus    sessionState = iota
	inputFocus    sessionState = iota
	dialogFocus   sessionState = iota
	ContainerMode              = ":containers"
	ImageMode                  = ":images"
	VolumeMode                 = ":volumes"
)

var (
	KeyQuit    = key.NewBinding(key.WithKeys("ctrl+c"))
	KeyCommand = key.NewBinding(key.WithKeys(":"))
	KeyEnter   = key.NewBinding(key.WithKeys("enter"))
	KeyEscape  = key.NewBinding(key.WithKeys("esc"))
	KeySpace   = key.NewBinding(key.WithKeys(" "))
	KeyRight   = key.NewBinding(key.WithKeys("right", "l"))
	KeyLeft    = key.NewBinding(key.WithKeys("left", "h"))
)

type model struct {
	state  sessionState
	banner banner.Model
	table  table.Model
	input  textinput.Model
	index  int
	width  int
	height int
	mode   string
	dialog dialog.Model
}

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

func (m model) Init() tea.Cmd {
	logger.Trace()

	return tea.Batch(
		m.banner.Init(),
		m.table.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Trace(msg)
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case dialogFocus:
			if m.dialog.ConfirmActions(msg) {
				m.state = tableFocus
				return m, nil
			}
		case inputFocus:
			switch {
			case key.Matches(msg, KeyEnter):
				m.state = tableFocus
				err := m.setContext(m.input.Value())
				if err != nil {
					m.state = dialogFocus
					m.dialog = dialog.NewDialog("ERROR", err.Error(), "Dismiss")
					log.Printf("Context Error: [%v]", err)
				}
				if strings.HasPrefix(":quit", m.input.Value()) {
					return m, tea.Quit
				}
				m.input.SetValue("")
				m.input.Prompt = ""
			}
		case tableFocus:
			switch {
			case key.Matches(msg, KeyQuit):
				return m, tea.Quit
			case key.Matches(msg, KeyCommand):
				if m.table.Focus() != table.DialogFocus {
					m.state = inputFocus
					m.input.Prompt = "!"
					m.input.SetValue("")
					m.input.Focus()
				}
			case key.Matches(msg, KeyEscape):
				err := m.setContext(m.mode)
				if err != nil {
					log.Printf("Context Error: [%v]", err)
				}
			}
		}

		// route keymsg to the correct widget
		switch m.state {
		case tableFocus:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		case inputFocus:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	case timers.TimerMsg:
		m.banner, cmd = m.banner.Update(msg)
		cmds = append(cmds, cmd)
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func (m *model) setContext(name string) error {
	logger.Trace(name)
	switch {
	case strings.HasPrefix(ContainerMode, name):
		m.table.SetContext(table.ContainerContext)
	case strings.HasPrefix(ImageMode, name):
		m.table.SetContext(table.ImageContext)
	case strings.HasPrefix(VolumeMode, name):
		m.table.SetContext(table.VolumeContext)
	default:
		return errors.New("Unsupported Command: [" + name + "]")
	}
	m.mode = name
	return nil
}

func (m model) View() string {
	logger.Trace()
	var s string

	banner := lipgloss.NewStyle().
		Render(m.banner.View())

	table := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Render(m.table.View())

	color := "22"
	if m.state == inputFocus {
		color = "69"
	}

	input := lipgloss.NewStyle().
		Width(m.width-2).
		Height(1).
		Align(lipgloss.Left, lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(color)).
		Render(m.input.View())

	s += lipgloss.JoinVertical(lipgloss.Top,
		banner,
		table,
		input,
	)

	if m.state == dialogFocus {
		popup := m.dialog.ConfirmDialog()

		return dialog.PlaceOverlay(
			lipgloss.Width(s)/2-lipgloss.Width(popup)/2,
			lipgloss.Height(s)/2-lipgloss.Height(popup)/2,
			popup,
			s,
			false,
		)
	}

	return s
}

func Show(mode string) {
	tea.LogToFile("debug.log", "")
	logger.Setup()
	logger.StartBanner()

	m := model{state: tableFocus}
	m.banner = banner.New()
	m.table = table.New()
	m.input = textinput.New()
	m.input.Prompt = ""
	m.setContext(mode)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
