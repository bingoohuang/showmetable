package util

import (
	"fmt"
	_ "net/http/pprof" // nolint gosec
	"os"

	"github.com/bingoohuang/gou/cnf"
	"github.com/bingoohuang/gou/htt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// InitFlags ...
func InitFlags() {
	help := pflag.BoolP("help", "h", false, "help")
	pflag.StringP("fmt", "f", "md", "md/excel")
	pflag.StringP("ds", "d", "", `eg. 
	MYSQL_PWD=8BE4 mysql -h 127.0.0.1 -P 9633 -u root
	mysql -h 127.0.0.1 -P 9633 -u root -p8BE4
	mysql -h 127.0.0.1 -P 9633 -u root -p8BE4 -Dtest
	mysql -h127.0.0.1 -u root -p8BE4 -Dtest
	127.0.0.1:9633 root/8BE4
	127.0.0.1 root/8BE4
	127.0.0.1:9633 root/8BE4 db=test
	root:8BE4@tcp(127.0.0.1:9633)/?charset=utf8mb4&parseTime=true&loc=Local
`)
	pflag.StringP("tables", "t", "", "show tables, separated by comma, empty for all")
	pflag.BoolP("verbose", "V", false, "verbose details")

	pprofAddr := htt.PprofAddrPflag()

	pflag.Parse()
	cnf.CheckUnknownPFlags()

	if *help {
		fmt.Printf("Built on %s from sha1 %s\n", Compile, Version)
		pflag.PrintDefaults()
		os.Exit(0)
	}

	htt.StartPprof(*pprofAddr)

	viper.SetEnvPrefix("SHOWMETABLE")
	viper.AutomaticEnv()

	_ = viper.BindPFlags(pflag.CommandLine)
}
