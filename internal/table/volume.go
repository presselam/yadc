package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/presselam/yadc/internal/bubble"
	"github.com/presselam/yadc/internal/dialog"
	"github.com/presselam/yadc/internal/docker"
	"github.com/presselam/yadc/internal/logger"
)

func (m *Model) volumeActions() []KeyMapping {
	retval := []KeyMapping{
		{cmd: (*Model).inspectVolume,
			key: key.NewBinding(
				key.WithKeys("i"),
				key.WithHelp("i", "inspect"),
			),
		},
		{cmd: (*Model).removeVolume,
			key: key.NewBinding(
				key.WithKeys("ctrl+d"),
				key.WithHelp("ctrl+d", "remove"),
			),
		},
		{cmd: (*Model).pruneVolumes,
			key: key.NewBinding(
				key.WithKeys("ctrl+p"),
				key.WithHelp("ctrl+p", "prune"),
			),
		},
	}

	return retval
}

func (m *Model) PopulateVolumes() error {
	logger.Trace()
	results, err := docker.Volumes()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	columns := []bubble.Column{}
	for i, col := range results.Columns {
		columns = append(columns, bubble.Column{Title: col, Width: results.Width[i]})
	}

	rows := []bubble.Row{}
	for _, r := range results.Data {
		rows = append(rows, r)
	}

	m.table.SetData(columns, rows)
	m.sortRows()
	return nil
}

func (m *Model) removeVolume(id string) {
	logger.Trace()
	logger.Debug("here")
	_, err := docker.VolumeDelete(id)
	if err != nil {
		logger.Error(err)
	}
	m.PopulateVolumes()
}

func (m *Model) pruneVolumes(id string) {
	logger.Trace(id)
	if m.focus == TableFocus {
		m.focus = DialogFocus
		m.confirm = dialog.NewDialog(
			"Prune",
			"This will remove all unused volumes",
			"Confirm", "Dismiss",
		)
	} else {
		go docker.VolumesPrune(id)
		m.PopulateVolumes()
	}
}
