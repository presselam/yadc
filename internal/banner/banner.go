package banner

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/presselam/yadc/internal/docker"
	"github.com/presselam/yadc/internal/timers"
	"log"
	"strconv"
	"time"
)

type Model struct {
	id   int
	info docker.ServerInfo
}

var titleStyle = lipgloss.NewStyle().
	Bold(false).
	Foreground(lipgloss.Color("70"))

var valueStyle = lipgloss.NewStyle().
	Bold(true).
	AlignHorizontal(lipgloss.Right).
	Foreground(lipgloss.Color("255"))

func (m Model) tick() tea.Cmd {
	delay := 5 * time.Second

	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return timers.TimerMsg{ID: m.id, Tag: t, Timeout: false}
	})
}

func (m Model) Init() tea.Cmd {
	log.Println("table.init")
	return m.tick()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	log.Printf("banner.update: [%v]", msg)

	switch msg := msg.(type) {
	case timers.TimerMsg:
		if msg.ID == m.id {
			status, err := docker.Info()
			if err != nil {
				log.Println(err)
			}
			m.info = status
		}
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	}

	return m, m.tick()
}

func (m Model) View() string {
	var s string
	s += lipgloss.JoinVertical(lipgloss.Top,
		displayField("Server:    ", m.info.Name),
		displayField("Server Ver:", m.info.ServerVersion),
		displayField("Client Ver:", m.info.ClientVersion),
		displayField("Images:    ", strconv.Itoa(m.info.Images)),
		displayField("Containers:", fmt.Sprintf("%d / %d / %d", m.info.Running, m.info.Paused, m.info.Stopped)),
	)

	return s
}

func displayField(name string, value string) string {
	var s string
	s = fmt.Sprintf("%s %s",
		titleStyle.Render(name),
		valueStyle.Render(value),
	)

	return s
}

func (m *Model) resize(width int, height int) {
}

func New() Model {
	status, err := docker.Info()
	if err != nil {
		log.Println(err)
	}

	m := Model{
		id:   timers.NextID(),
		info: status,
	}
	return m
}
