module github.com/GoDannyLai/binlog_rollback

go 1.25.9

replace github.com/siddontang/go-mysql => ./vendor/github.com/siddontang/go-mysql

replace github.com/dropbox/godropbox => ./vendor/github.com/dropbox/godropbox

replace github.com/jingchengli/dannytools => ./local/dannytools

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/dropbox/godropbox v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.9.3
	github.com/jingchengli/dannytools v0.0.0-20180809154344-1107578f18ef
	github.com/juju/errors v1.0.0
	github.com/siddontang/go-mysql v0.0.0-00010101000000-000000000000
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/toolkits/slice v0.0.0-20141116085117-e44a80af2484
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/pingcap/check v0.0.0-20211026125417-57bd13f7b5f0 // indirect
	github.com/pingcap/errors v0.11.4 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726 // indirect
	github.com/siddontang/go-log v0.0.0-20190221022429-1e957dd83bed // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	golang.org/x/sys v0.13.0 // indirect
)
