package main

import (
	"fmt"
	"path/filepath"

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

	for _, st := range showTableModels {
		t := st.GetScheme() + "." + st.GetName()
		columns, err := a.schema.TableColumns(t)

		if err != nil {
			panic(err)
		}

		fmt.Printf("table %s, columns: %+v\n", t, columns)
	}
}

func findTable(showTables []string, t model.Table) model.Table {
	for _, st := range showTables {
		if yes, _ := filepath.Match(st, t.GetName()); yes {
			return t
		}
	}

	return nil
}
