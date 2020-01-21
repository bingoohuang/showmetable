package main

import (
	"path/filepath"
	"regexp"
	"time"

	"github.com/bingoohuang/gou/str"

	"github.com/bingoohuang/showmetable/dbi"
	"github.com/bingoohuang/showmetable/model"
	"github.com/bingoohuang/showmetable/util"
	"github.com/spf13/viper"
)

func main() {
	util.InitFlags()

	app := MakeApp()
	app.showTables()
}

// App ...
type App struct {
	schema model.DbSchema
}

// MakeApp ...
func MakeApp() *App {
	schema, err := dbi.CreateMySQLSchema(viper.GetString("ds"))
	if err != nil {
		panic(err)
	}

	schema.SetVerbose(viper.GetBool("verbose"))

	return &App{schema: schema}
}

func (a *App) showTables() {
	tables, err := a.schema.Tables()
	if err != nil {
		panic(err)
	}

	showTables := viper.GetStringSlice("tables")
	showAll := len(showTables) == 0

	showTableModels := make([]model.Table, 0)

	for _, t := range tables {
		if showAll {
			showTableModels = append(showTableModels, t)
		} else if mt := findTable(showTables, t); mt != nil {
			showTableModels = append(showTableModels, mt)
		}
	}

	w := MakeTableWriters(makeMarkdownWriter(), makeExcelWriter())

	for _, st := range showTableModels {
		a.showTable(st, w)
	}

	_, _ = w.SaveAs(time.Now().Format("20060102150405"))
}

var newlineRe = regexp.MustCompile(`\r?\n`) // nolint gochecknoglobals

func (a *App) showTable(st model.Table, w TableWriter) {
	t := st.GetScheme() + "." + st.GetName()
	columns, err := a.schema.TableColumns(t)

	if err != nil {
		panic(err)
	}

	w.WriteTable(st)
	w.WriteColumns(columns)
}

func findTable(showTables []string, t model.Table) model.Table {
	for _, st := range showTables {
		schema, tableName := str.Split2(st, ".", true, true)
		if !schemaMatch(schema, t.GetScheme()) {
			continue
		}

		if yes, _ := filepath.Match(tableName, t.GetName()); yes {
			return t
		}
	}

	return nil
}

func schemaMatch(showSchema, tableSchema string) bool {
	return showSchema == "" || showSchema == tableSchema
}
