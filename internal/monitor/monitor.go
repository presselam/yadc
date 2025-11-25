package monitor

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/presselam/yadc/internal/banner"
	"github.com/presselam/yadc/internal/table"
	"github.com/presselam/yadc/internal/timers"
	"log"
	"strings"
	"time"
)

type sessionState uint

const (
	defaultTime              = time.Minute
	tableFocus  sessionState = iota
	inputFocus
	containerMode = ":containers"
	imageMode     = ":images"
	valueMode     = ":volumes"
)

var (
	KeyQuit    = key.NewBinding(key.WithKeys("ctrl+c"))
	KeyCommand = key.NewBinding(key.WithKeys(":"))
	KeyEnter   = key.NewBinding(key.WithKeys("enter"))
	KeyEscape  = key.NewBinding(key.WithKeys("esc"))
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
}

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.banner.Init(),
		m.table.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("monitor.update: [%v]", msg)
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyQuit):
			return m, tea.Quit
		case key.Matches(msg, KeyCommand):
			m.state = inputFocus
			m.input.Prompt = "!"
			m.input.SetValue("")
			m.input.Focus()
		case key.Matches(msg, KeyEnter):
			m.state = tableFocus
			err := m.setContext(m.input.Value())
			if err != nil {
				log.Printf("Context Error: [%v]", err)
			}
			if strings.HasPrefix(":quit", m.input.Value()) {
				return m, tea.Quit
			}
			m.input.SetValue("")
			m.input.Prompt = ""
		case key.Matches(msg, KeyEscape):
			err := m.setContext(m.mode)
			if err != nil {
				log.Printf("Context Error: [%v]", err)
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
	switch {
	case strings.HasPrefix(containerMode, name):
		m.table.SetContext(table.ContainerContext)
	case strings.HasPrefix(imageMode, name):
		m.table.SetContext(table.ImageContext)
	case strings.HasPrefix(valueMode, name):
		m.table.SetContext(table.VolumeContext)
	default:
		return errors.New("Unsupported Command: [" + name + "]")
	}
	m.mode = name
	return nil
}

func (m model) View() string {
	var s string

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

	return s
}

func Show() {
	tea.LogToFile("debug.log", "debug")
	log.Println("#=======================================================")
	log.Println("# ", time.Now())
	log.Println("#=======================================================")

	m := model{state: tableFocus, mode: imageMode}
	m.banner = banner.New()
	m.table = table.New()
	m.input = textinput.New()
	m.input.Prompt = ""

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
