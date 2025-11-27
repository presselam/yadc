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
	"strconv"
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
	modal  string
	button bool
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
		switch {
		case key.Matches(msg, KeyRight):
			if m.state == dialogFocus {
				m.button = !m.button
				return m, nil
			}
		case key.Matches(msg, KeyLeft):
			if m.state == dialogFocus {
				m.button = !m.button
				return m, nil
			}
		case key.Matches(msg, KeySpace):
			if m.state == dialogFocus {
				m.state = tableFocus
				m.modal = ""
				return m, nil
			}
		case key.Matches(msg, KeyQuit):
			return m, tea.Quit
		case key.Matches(msg, KeyCommand):
			m.state = inputFocus
			m.input.Prompt = "!"
			m.input.SetValue("")
			m.input.Focus()
		case key.Matches(msg, KeyEnter):
			if m.state == dialogFocus {
				m.state = tableFocus
				m.modal = ""
				logger.Info("User selected: [", strconv.FormatBool(m.button), "]")
				return m, nil
			} else {
				m.state = tableFocus
				err := m.setContext(m.input.Value())
				if err != nil {
					m.state = dialogFocus
					m.modal = err.Error()
					m.button = true
					log.Printf("Context Error: [%v]", err)
				}
				if strings.HasPrefix(":quit", m.input.Value()) {
					return m, tea.Quit
				}
				m.input.SetValue("")
				m.input.Prompt = ""
			}
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
		dialogBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

		buttonStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginRight(2).
			MarginTop(1)

		activeButtonStyle := buttonStyle.
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#F25D94")).
			MarginRight(2).
			Underline(true)

		var okButton string
		var cancelButton string
		if m.button {
			okButton = activeButtonStyle.Render("Yes")
			cancelButton = buttonStyle.Render("Maybe")
		} else {
			okButton = buttonStyle.Render("Yes")
			cancelButton = activeButtonStyle.Render("Maybe")
		}

		//		blends := gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 50)

		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(m.modal)
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

		confirm := lipgloss.Place(m.width-2, 9,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
		)
		return dialog.PlaceOverlay(
			lipgloss.Width(table)/2,
			lipgloss.Height(table)/2,
			confirm,
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
