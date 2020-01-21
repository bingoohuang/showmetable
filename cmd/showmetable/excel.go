package main

import (
	"bytes"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/showmetable/model"
	"github.com/bingoohuang/statiq/fs"

	// import static resources
	_ "github.com/bingoohuang/showmetable/statiq"
)

type excelWriter struct {
	lineIndex int
	f         *excelize.File
}

const sheet1 = "Sheet1"

// WriteTable writes table information
func (w *excelWriter) WriteTable(st model.Table) {
	a := fmt.Sprintf("A%d", w.lineIndex)
	b := fmt.Sprintf("B%d", w.lineIndex)
	c := fmt.Sprintf("C%d", w.lineIndex)
	d := fmt.Sprintf("D%d", w.lineIndex)
	f := fmt.Sprintf("F%d", w.lineIndex)

	_ = w.f.SetCellValue(sheet1, a, "表名")
	_ = w.f.SetCellStyle(sheet1, a, a, w.getCellStyle("A1"))

	_ = w.f.MergeCell(sheet1, b, c)
	_ = w.f.SetCellValue(sheet1, b, st.GetName())
	_ = w.f.SetCellStyle(sheet1, b, b, w.getCellStyle("B1"))

	_ = w.f.MergeCell(sheet1, d, f)
	_ = w.f.SetCellValue(sheet1, d, st.GetComment())
	_ = w.f.SetCellStyle(sheet1, d, d, w.getCellStyle("D1"))

	w.lineIndex += 2
}

// WriteColumns writes columns of table
func (w *excelWriter) WriteColumns(columns []model.TableColumn) {
	w.writeHeadCells("A", "#", "字段名称", "数据类型", "是否可空", "默认值", "字段注释")

	w.lineIndex++

	for i, col := range columns {
		seq := fmt.Sprintf("%d", i+1) // nolint gomnd
		w.writeDataCells(i, "A", seq, col.GetName(), col.GetDataType(),
			str.If(col.IsNullable(), "Y", "N"), col.GetDefault(), col.GetComment())

		w.lineIndex++
	}

	w.lineIndex++
}

// SaveAs save as filename
func (w *excelWriter) SaveAs(filename string) (string, error) {
	f := filename + ".xlsx"
	return f, w.f.SaveAs(f)
}

func (w *excelWriter) writeHeadCells(colStart string, values ...string) {
	for i, v := range values {
		col := columnOffset(colStart, i)
		w.writeCell(col, v, col+"3")
	}
}

func (w *excelWriter) writeDataCells(seq int, colStart string, values ...string) {
	for i, v := range values {
		col := columnOffset(colStart, i)

		w.writeCell(col, v, col+str.If(seq%2 == 0, "4", "5"))
	}
}

func columnOffset(colStart string, i int) string {
	return string(int(colStart[0]) + i)
}

func (w *excelWriter) writeCell(columnAxis, value, styleCellAxis string) {
	a := fmt.Sprintf("%s%d", columnAxis, w.lineIndex)
	_ = w.f.SetCellValue(sheet1, a, value)

	style := w.getCellStyle(styleCellAxis)

	_ = w.f.SetCellStyle(sheet1, a, a, style)
}

func (w *excelWriter) getCellStyle(axis string) int {
	style, _ := w.f.GetCellStyle(sheet1, axis)

	return style
}

func makeExcelWriter() *excelWriter {
	tmplFs, _ := fs.New()
	tmplXlsx := bytes.NewReader(tmplFs.Files["/tmpl.xlsx"].Data)
	tmpl, _ := excelize.OpenReader(tmplXlsx)

	_ = tmpl.SetSheetPrOptions(sheet1, excelize.FitToPage(true))
	_ = tmpl.SetPageLayout(sheet1, excelize.FitToHeight(2), excelize.FitToWidth(2)) // nolint gomnd

	return &excelWriter{
		lineIndex: 1, // nolint gomnd
		f:         tmpl,
	}
}
