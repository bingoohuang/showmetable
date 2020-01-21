package main

import (
	"fmt"

	"github.com/bingoohuang/showmetable/model"
)

// TableWriter write table definitions.
type TableWriter interface {
	// WriteTable writes table information
	WriteTable(st model.Table)
	// WriteColumns writes columns of table
	WriteColumns(columns []model.TableColumn)
	// SaveAs saves to filename without extension
	SaveAs(filename string) (string, error)
}

// TableWriters composes multiple TableWriter
type TableWriters struct {
	writers []TableWriter
}

// MakeTableWriters makes a composed TableWriter
func MakeTableWriters(writers ...TableWriter) TableWriter {
	return &TableWriters{writers: writers}
}

// WriteTable writes table information
func (w *TableWriters) WriteTable(st model.Table) {
	for _, s := range w.writers {
		s.WriteTable(st)
	}
}

// WriteColumns writes columns of table
func (w *TableWriters) WriteColumns(columns []model.TableColumn) {
	for _, s := range w.writers {
		s.WriteColumns(columns)
	}
}

// SaveAs saves to filename without extension
func (w *TableWriters) SaveAs(filename string) (string, error) {
	for _, s := range w.writers {
		if f, err := s.SaveAs(filename); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(f, "generated!")
		}
	}

	return "", nil
}
