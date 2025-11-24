package table

import (
	"github.com/presselam/yadc/internal/bubble"
	"github.com/presselam/yadc/internal/docker"
	"log"
	// "github.com/charmbracelet/bubbles/key"
)

func (m *Model) inspectContainer(id string) {
	log.Printf("container.inspect.%s", id)
	m.context = InspectContext
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
