package main

import (
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
	schema.SetShowTables(viper.GetStringSlice("tables"))

	return &App{schema: schema}
}

func (a *App) showTables() {

}
