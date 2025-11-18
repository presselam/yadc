package table

import (
  "github.com/presselam/yadc/internal/docker"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
  BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	 table table.Model
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
				case "q", "ctrl+c": 
				  return m, tea.Quit
				case "enter":
				  return m, tea.Batch(
					  tea.Printf("selected: [%s]", m.table.SelectedRow()[1]),
					)	
			}
	}

  m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func (m Model) resize(width int, height int) {

  m.table.SetWidth(width)
}

func New() Model{
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
	}

  containers, _ := docker.Containers()
  rows := []table.Row{}
  for _, c := range containers {
    rows = append(rows, table.Row{c.ID, c.Names[0], c.Image, c.State})
  }



	t := table.New(
	  table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)	

	s := table.DefaultStyles()
	s.Header = s.Header.
	  BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
	  Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	m := Model{t}
	
	return m
}
