package table

import (
  "log"
  "github.com/presselam/yadc/internal/docker"
  "github.com/presselam/yadc/internal/bubble"
)

func (m *Model) PopulateImages() error {
  results, err := docker.Images()
  if err != nil {
    return err
  }

  columns := []bubble.Column{}
  for i, col := range results.Columns {
    log.Println("Column: ["+ col + "]")
    columns = append(columns, bubble.Column{Title: col, Width: results.Width[i]})
  }

  rows := []bubble.Row{}
  for _, r := range results.Data {
    log.Println(r)
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

func (m *Model) imageHandler(command string) {
  switch command {
    case "ctrl+d":
      row := m.table.SelectedRow()
      log.Printf("deleting image: [%s]\n", row[0]) 
      _, err := docker.ImageDelete(row[0])
      if err != nil {
        log.Printf("error deleting image: [%v]\n", err) 
      }
  }
}
