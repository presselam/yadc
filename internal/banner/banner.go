package banner

import (
  "fmt"
  "log"
  "github.com/presselam/yadc/internal/docker"
   tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
  "strconv"
)

type Model struct {
  info docker.ServerInfo 
}

var titleStyle = lipgloss.NewStyle().
  Bold(false).
  Foreground(lipgloss.Color("70"))

var valueStyle = lipgloss.NewStyle().
  Bold(true).
  AlignHorizontal(lipgloss.Right).
  Foreground(lipgloss.Color("255"))



func (m Model) Init() tea.Cmd {
  return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

  switch msg := msg.(type) {
    case tea.WindowSizeMsg:
      m.resize(msg.Width, msg.Height)
  }

  return m, nil
}

func (m Model) View() string {
  var s string
   s += lipgloss.JoinVertical(lipgloss.Top,
      displayField("Server:    ", m.info.Name),
      displayField("Server Ver:", m.info.ServerVersion),
      displayField("Client Ver:", m.info.ClientVersion),
      displayField("Images:    ", strconv.Itoa( m.info.Images)),
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
  status,err := docker.Info()
  if err != nil {
    log.Println(err)
  }

  m := Model{status}
  return m 
}

