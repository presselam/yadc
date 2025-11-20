package table

import (
  "log"
  "github.com/presselam/yadc/internal/bubble"
  "github.com/presselam/yadc/internal/docker"
)

func (m *Model) PopulateContainers() error {
  results, err := docker.Containers()
  if err != nil{
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

  if len(m.table.Columns()) < len(columns) {
    m.table.SetColumns(columns)
    m.table.SetRows(rows)
  }else{
    m.table.SetRows(rows)
    m.table.SetColumns(columns)
  }   

  return nil
}

func (m *Model) containerHandler(command string) {
  switch command {
    case "ctrl+d":
      row := m.table.SelectedRow()
      log.Printf("deleting image: [%s]\n", row[0]) 
  }
}
