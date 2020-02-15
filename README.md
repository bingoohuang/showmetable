# showmetable

show me the markdown/excel tables from database tables

## build

`go fmt ./...&&goimports -w .&&golint ./...&&golangci-lint run --enable-all&& go install  -ldflags="-s -w" ./...`

## run

```bash
$ showmetable -h                                                                         
Built on 2020-01-21 14:37:35 +0800 by go version go1.13.5 darwin/amd64 from sha1 2020-01-21 11:32:37 +0800 @f77551db719c1f0bf2facecdb40fd5faaf4827f6 @
  -d, --ds string           eg. 
                                MYSQL_PWD=8BE4 mysql -h 127.0.0.1 -P 9633 -u root
                                mysql -h 127.0.0.1 -P 9633 -u root -p8BE4
                                mysql -h 127.0.0.1 -P 9633 -u root -p8BE4 -Dtest
                                mysql -h127.0.0.1 -u root -p8BE4 -Dtest
                                127.0.0.1:9633 root/8BE4
                                127.0.0.1 root/8BE4
                                127.0.0.1:9633 root/8BE4 db=test
                                root:8BE4@tcp(127.0.0.1:9633)/?charset=utf8mb4&parseTime=true&loc=Local
                            
  -h, --help                help
      --pprof-addr string   pprof address to listen on, not activate pprof if empty, eg. --pprof-addr localhost:6060
  -t, --tables string       show tables, separated by comma, empty for all, like db1.*,db1.mytable
  -V, --verbose             verbose details

$ showmetable -d "1.2.1.1:3306 aaa/bbb"  -t "api_manager.alarm_*"            
INFO[0000] dataSourceName:aaa:bbb@tcp(1.2.1.1:3306)/?charset=utf8mb4&parseTime=true&loc=Local 
20200121144038.md generated!
20200121144038.xlsx generated!

```

## TODO

1. support direct ddl sql file, like [Metaprogram in Go How to use Golang to metaprogram a CRUD API](https://levelup.gitconnected.com/metaprogram-in-go-5a2a7e989613) and [code](https://github.com/exyzzy/metaapi)
