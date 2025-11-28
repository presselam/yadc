package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/presselam/yadc/internal/bubble"
	"github.com/presselam/yadc/internal/dialog"
	"github.com/presselam/yadc/internal/docker"
	"github.com/presselam/yadc/internal/logger"
)

func (m *Model) PopulateContainers() error {
	results, err := docker.Containers()
	if err != nil {
		return err
	}

	total := 0
	columns := []bubble.Column{}
	for i, col := range results.Columns {
		columns = append(columns, bubble.Column{Title: col, Width: results.Width[i]})
		total += results.Width[i]
	}

	rows := []bubble.Row{}
	for _, r := range results.Data {
		rows = append(rows, r)
	}

	m.table.SetData(columns, rows)
	m.sortRows()
	return nil
}

func (m *Model) containerKeyMapping() []KeyMapping {
	retval := []KeyMapping{
		{cmd: (*Model).inspectContainer,
			key: key.NewBinding(
				key.WithKeys("i"),
				key.WithHelp("i", "inspect"),
			),
		},
		{cmd: (*Model).restartContainer,
			key: key.NewBinding(
				key.WithKeys("ctrl+r", "ctrl+s"),
				key.WithHelp("ctrl+r", "restart"),
			),
		},
		{cmd: (*Model).stopContainer,
			key: key.NewBinding(
				key.WithKeys("ctrl+k"),
				key.WithHelp("ctrl+s", "start"),
			),
		},
		{cmd: (*Model).pruneContainer,
			key: key.NewBinding(
				key.WithKeys("ctrl+p"),
				key.WithHelp("ctrl+p", "prune"),
			),
		},
		{cmd: (*Model).logContainer,
			key: key.NewBinding(
				key.WithKeys("l"),
				key.WithHelp("l", "logs"),
			),
		},
	}

	return retval
}

func (m *Model) restartContainer(id string) {
	go docker.ContainerRestart(id)
	m.PopulateContainers()
}

func (m *Model) stopContainer(id string) {
	go docker.ContainerStop(id)
	m.PopulateContainers()
}

func (m *Model) pruneContainer(id string) {
	logger.Debug("table.containers.pruneContainer")
	if m.focus == tableFocus {
		m.focus = dialogFocus
		m.confirm = dialog.NewDialog(
			"Prune",
			"This will remove all stopped containers",
			"Confirm", "Dismiss",
		)
	}

	go docker.ContainerPrune(id)
	m.PopulateContainers()
}

func (m *Model) logContainer(id string) {
	m.selected = id
	m.SetContext(LogsContext)
	m.FetchLogs()
}

func (m *Model) FetchLogs() error {
	results, err := docker.ContainerLog(m.selected, "15m")
	if err != nil {
		return err
	}

	total := 0
	columns := []bubble.Column{}
	for i, col := range results.Columns {
		columns = append(columns, bubble.Column{Title: col, Width: results.Width[i]})
		total += results.Width[i]
	}

	rows := []bubble.Row{}
	for _, r := range results.Data {
		rows = append(rows, r)
	}

	m.table.SetData(columns, rows)
	m.table.SetCursor(len(rows))

	return nil
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
