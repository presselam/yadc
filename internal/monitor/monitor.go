package monitor

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/presselam/yadc/internal/banner"
	"github.com/presselam/yadc/internal/table"
	"log"
	"strings"
	"time"
)

type sessionState uint

const (
	defaultTime              = time.Minute
	tableFocus  sessionState = iota
	inputFocus
)

type model struct {
	state   sessionState
	banner  banner.Model
	table   table.Model
	input   textinput.Model
	index   int
	width   int
	height  int
	command string
}

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case ":":
			m.state = inputFocus
			m.input.Prompt = "!"
			m.input.SetValue("")
			m.input.Focus()
		case "enter":
			m.state = tableFocus
			if m.input.Value() == ":q" {
				return m, tea.Quit
			}
			err := m.setContext(m.input.Value())
			if err != nil {
				log.Printf("Context Error: [%v]", err)
			}

			m.command = m.input.Value()
			m.input.SetValue("")
			m.input.Prompt = ""
		}
		switch m.state {
		case tableFocus:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		case inputFocus:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.input, cmd = m.input.Update(msg)
		m.table, cmd = m.table.Update(msg)
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func (m *model) setContext(name string) error {
	switch {
	case strings.HasPrefix(":containers", name):
		m.table.SetContext(table.ContainerContext)
	case strings.HasPrefix(":images", name):
		m.table.SetContext(table.ImageContext)
	case strings.HasPrefix(":volumes", name):
		m.table.SetContext(table.VolumeContext)
	default:
		return errors.New("Unsupported Command: [" + name + "]")
	}
	return nil
}

func (m *model) currentFocusedModel() string {
	if m.state == inputFocus {
		return "input"
	}
	return "table"
}

func (m model) View() string {
	var s string
	model := m.currentFocusedModel()

	inputStyle := lipgloss.NewStyle().
		Width(m.width-2).
		Height(1).
		Align(lipgloss.Left, lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("22"))

	tableStyle := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center)

		//  focusedStyle := lipgloss.NewStyle().
		//  	Width(m.width-2).
		//		Height(15).
		//		Align(lipgloss.Center, lipgloss.Center).
		//		BorderStyle(lipgloss.NormalBorder()).
		//		BorderForeground(lipgloss.Color("69"))

	if m.state == tableFocus {
		s += lipgloss.JoinVertical(lipgloss.Top,
			fmt.Sprintf("%s", m.banner.View()),
			tableStyle.Render(m.table.View()),
			inputStyle.Render(m.input.View()),
		)
	} else if m.state == inputFocus {
		inputStyle = inputStyle.BorderForeground(lipgloss.Color("69"))
		s += lipgloss.JoinVertical(lipgloss.Top,
			fmt.Sprintf("%s", m.banner.View()),
			tableStyle.Render(m.table.View()),
			inputStyle.Render(m.input.View()),
		)
	} else {
		s += lipgloss.JoinVertical(lipgloss.Top,
			fmt.Sprintf("%s", m.banner.View()),
			tableStyle.Render(m.table.View()),
			inputStyle.Render(m.input.View()),
		)
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next * n: new %s q: exit\n", model))

	return s
}

func Show() {
	tea.LogToFile("debug.log", "debug")
	m := model{state: tableFocus}
	m.banner = banner.New()
	m.table = table.New()
	m.input = textinput.New()
	m.input.Prompt = ""

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
