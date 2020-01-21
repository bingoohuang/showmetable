package main

import (
	"io/ioutil"
	"strings"

	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/showmetable/model"
)

type markdownWriter struct {
	md strings.Builder
}

func makeMarkdownWriter() *markdownWriter {
	return &markdownWriter{}
}

func (w *markdownWriter) write(ss ...string) {
	for _, s := range ss {
		_, _ = w.md.WriteString(s)
	}
}

// SaveAs save as filename
func (w *markdownWriter) SaveAs(filename string) (string, error) {
	f := filename + ".md"
	return f, ioutil.WriteFile(f, []byte(w.md.String()), 0644)
}

// WriteTable writes table information
func (w *markdownWriter) WriteTable(st model.Table) {
	w.write("## ", st.GetName(), "\n\n")

	comment := st.GetComment()
	if comment != "" {
		w.write(comment, "\n\n")
	}
}

// WriteColumns writes columns of table
func (w *markdownWriter) WriteColumns(columns []model.TableColumn) {
	w.write("字段名称|数据类型|是否可空|默认值|字段注释\n")
	w.write("---|  ---    |  ---   |  ---  |  ---\n")

	for _, c := range columns {
		cc := c.GetComment()
		if cc != "" {
			cc = newlineRe.ReplaceAllString(cc, "<br>")
		}

		w.write(c.GetName(), "|", c.GetDataType(), "|", str.If(c.IsNullable(),
			"Y", "N"), "|", c.GetDefault(), "|", cc, "\n")
	}

	w.write("\n")
}
