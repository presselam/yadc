package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	"github.com/presselam/yadc/internal/bubble"
)

func ContainerFormatter(row bubble.Row) lipgloss.Style {
	style := lipgloss.NewStyle()

	if len(row) > 4 {
		switch row[3] {
		case container.StateCreated:
			style = style.Foreground(lipgloss.Color("82"))
		case container.StatePaused:
			style = style.Foreground(lipgloss.Color("64"))
		case container.StateRestarting:
			style = style.Foreground(lipgloss.Color("202"))
		case container.StateRemoving:
			style = style.Foreground(lipgloss.Color("169"))
		case container.StateExited:
			style = style.Foreground(lipgloss.Color("242"))
		case container.StateDead:
			style = style.Foreground(lipgloss.Color("196"))
			//case container.StateRunning:
		default:
			style = style.Foreground(lipgloss.Color("225"))
		}
	}

	return style
}

func ImageFormatter(row bubble.Row) lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("225"))

	if row[2] == "0" {
		style = style.Foreground(lipgloss.Color("242"))
	}
	if row[1] == "<none>" {
		style = style.Foreground(lipgloss.Color("196"))
	}

	return style
}
