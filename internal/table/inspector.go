package table

import (
	"github.com/presselam/yadc/internal/bubble"
	"github.com/presselam/yadc/internal/docker"
	"github.com/presselam/yadc/internal/logger"
)

func (m *Model) inspectContainer(id string) {
	logger.Debug("table.inspector.container.inspect.", id)
	m.SetContext(InspectContext)
	results, _ := docker.ContainerInspect(id)

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

func (m *Model) inspectImage(id string) {
	logger.Debug("tabel.inspector.image.inspect.", id)
	m.SetContext(InspectContext)
	results, _ := docker.ImageInspect(id)

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

func (m *Model) inspectVolume(id string) {
	logger.Debug("tabel.inspector.volume.inspect.", id)
	m.SetContext(InspectContext)
	results, _ := docker.VolumeInspect(id)

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
