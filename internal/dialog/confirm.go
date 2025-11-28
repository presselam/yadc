package dialog

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/presselam/yadc/internal/logger"
)

func (m *Model) ConfirmActions(msg tea.KeyMsg) bool {
	logger.Debug("dialog.confirm.actions:", msg.String())
	switch {
	case key.Matches(msg, KeyDismiss):
		m.confirmation = true
		return true
	case key.Matches(msg, KeyRight):
		m.selected++
		m.selected %= 2
		return true
	case key.Matches(msg, KeyLeft):
		m.selected--
		m.selected = max(m.selected, -m.selected) % 2
		return true
	}
	return false
}

func (m Model) ConfirmDialog() string {
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

	buttons := []string{}
	for i, lbl := range m.buttons {
		var s lipgloss.Style
		if i == m.selected {
			s = activeButtonStyle
		} else {
			s = buttonStyle
		}
		buttons = append(buttons, s.Render(lbl))
	}

	title := lipgloss.NewStyle().Width(50).Bold(true).Align(lipgloss.Center).Render(m.title)
	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(m.message)
	buttonBar := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)
	ui := lipgloss.JoinVertical(lipgloss.Center, title, question, buttonBar)

	return dialogBoxStyle.Render(ui)
}
