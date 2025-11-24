package table

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/presselam/yadc/internal/bubble"
	"github.com/presselam/yadc/internal/docker"
	"log"
)

func (m *Model) PopulateImages() error {
	results, err := docker.Images()
	if err != nil {
		return err
	}

	columns := []bubble.Column{}
	for i, col := range results.Columns {
		log.Println("Column: [" + col + "]")
		columns = append(columns, bubble.Column{Title: col, Width: results.Width[i]})
	}

	rows := []bubble.Row{}
	for _, r := range results.Data {
		log.Println(r)
		rows = append(rows, r)
	}

	m.table.SetData(columns, rows)
	return nil
}

func (m *Model) imageHandler(msg tea.KeyMsg) {
	switch msg.String() {
	case "ctrl+d":
		row := m.table.SelectedRow()
		log.Printf("deleting image: [%s]\n", row[0])
		_, err := docker.ImageDelete(row[0])
		if err != nil {
			log.Printf("error deleting image: [%v]\n", err)
		}
	}
}
