package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
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

	var md strings.Builder

	w := func(s string) {
		_, _ = md.WriteString(s)
	}

	for _, st := range showTableModels {
		a.showTable(st, w)
	}

	f := time.Now().Format("20060102150405") + ".md"
	if err := ioutil.WriteFile(f, []byte(md.String()), 0644); err != nil {
		panic(err)
	}

	fmt.Println(f, "generated!")
}

var newlineRe = regexp.MustCompile(`\r?\n`) // nolint gochecknoglobals

func (a *App) showTable(st model.Table, mdw func(s string)) {
	t := st.GetScheme() + "." + st.GetName()
	columns, err := a.schema.TableColumns(t)

	if err != nil {
		panic(err)
	}

	//fmt.Printf("table %s, columns: %+v\n", t, columns)

	createMarkdown(mdw, st, columns)
}

func createMarkdown(mdw func(s string), st model.Table, columns []model.TableColumn) {
	mdw("## ")
	mdw(st.GetName())
	mdw("\n\n")

	comment := st.GetComment()
	if comment != "" {
		mdw(comment)
		mdw("\n\n")
	}

	mdw("Name|DataType|Nullable|Default|Comment\n")
	mdw("---|  ---    |  ---   |  ---  |  ---\n")

	for _, c := range columns {
		cc := c.GetComment()
		if cc != "" {
			cc = newlineRe.ReplaceAllString(cc, "<br>")
		}

		mdw(c.GetName() +
			"|" + c.GetDataType() +
			"|" + str.If(c.IsNullable(), "Y", "N") +
			"|" + c.GetDefault() +
			"|" + cc + "\n")
	}

	mdw("\n")
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
