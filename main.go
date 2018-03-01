package main

import (
	"log"
	"scanproxy/scanproxy"

	"github.com/JimYJ/easysql/mysql"
)

var (
	dbhost       = ""
	dbport       = 3306
	dbname       = ""
	dbuser       = ""
	dbpass       = ""
	charset      = ""
	maxIdleConns = 500
	maxOpenConns = 500
)

func main() {
	initDBConn()
	scanproxy.InternetAllScan("CN")
}

func initDBConn() {
	mysql.Init(dbhost, dbport, dbname, dbuser, dbpass, charset, maxIdleConns, maxOpenConns)
	_, err := mysql.GetMysqlConn()
	if err != nil {
		log.Panic(err)
	}
}
