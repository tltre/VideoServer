package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	// 连接数据库应放在专门的配置文件中
	// DSN格式："username:password@tcp(host:port)/database?charset=xxx"
	dbConn, err = sql.Open("mysql", "root:sy20021213@tcp(localhost:9090)/videoserver?charset=utf8")
	if err != nil {
		panic(err)
	}
}
