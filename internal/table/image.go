package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/presselam/yadc/internal/bubble"
	"github.com/presselam/yadc/internal/dialog"
	"github.com/presselam/yadc/internal/docker"
	"github.com/presselam/yadc/internal/logger"
	"log"
)

func (m *Model) imageActions() []KeyMapping {
	retval := []KeyMapping{
		{cmd: (*Model).inspectImage,
			key: key.NewBinding(
				key.WithKeys("i"),
				key.WithHelp("i", "inspect"),
			),
		},
		{cmd: (*Model).historyImage,
			key: key.NewBinding(
				key.WithKeys("h"),
				key.WithHelp("ctrl+r", "restart"),
			),
		},
		{cmd: (*Model).removeImage,
			key: key.NewBinding(
				key.WithKeys("ctrl+d"),
				key.WithHelp("ctrl+d", "remove"),
			),
		},
		{cmd: (*Model).pruneImages,
			key: key.NewBinding(
				key.WithKeys("ctrl+p"),
				key.WithHelp("ctrl+p", "prune"),
			),
		},
		{cmd: (*Model).saveImage,
			key: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("l", "logs"),
			),
		},
	}

	return retval

}

func (m *Model) PopulateImages() error {
	results, err := docker.Images()
	if err != nil {
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

func (m *Model) historyImage(id string) {
	log.Printf("table.image.historyImage.%s", id)
	m.SetContext(InspectContext)
	results, _ := docker.ImageHistory(id)

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
}

func (m *Model) removeImage(id string) {
	log.Printf("table.image.removeImage.%s", id)
	go docker.ImageDelete(id)
	m.PopulateImages()
}

func (m *Model) pruneImages(id string) {
	logger.Trace(id)
	logger.Debug("table.image.pruneImage.%s", id)
	if m.focus == TableFocus {
		m.focus = DialogFocus
		m.confirm = dialog.NewDialog(
			"Prune",
			"This will remove all dangling images",
			"Confirm", "Dismiss",
		)
	} else {
		logger.Debug("table.image.pruneImage.confirmeds")
		go docker.ImagesPrune(id)
		m.PopulateImages()
	}
}

func (m *Model) saveImage(id string) {
	logger.Trace(id)
	logger.Debug("table.image.saveImage.%s", id)
	if m.focus == TableFocus {
		m.focus = DialogFocus
		m.confirm = dialog.NewDialog(
			"Save",
			"Create Image Tarball",
			"Confirm", "Dismiss",
		)
	} else {
		logger.Debug("table.image.saveImage.confirmed")
		go docker.ImageSave(id)
	}
}
