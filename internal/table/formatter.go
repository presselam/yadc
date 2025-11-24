package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/presselam/yadc/internal/bubble"
)

func ContainerFormatter(row bubble.Row) lipgloss.Style {
	style := lipgloss.NewStyle()

	if len(row) > 4 {
		switch row[3] {
		case "exited":
			style = style.Foreground(lipgloss.Color("247"))
		}
	}

	return style
}
