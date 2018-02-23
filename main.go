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
	// runtime.GOMAXPROCS(runtime.NumCPU())
	// ch := make(chan map[string]int, 1)
	// go scanproxy.CheckPort("192.168.10.242", 80, ch)
	// log.Println(<-ch)
	// iplist, total, totalPage, err := scanproxy.GetApnicIP("CN", 1, 1)
	// log.Println(iplist, total, totalPage, err)
	// allPortOk := scanproxy.ScanAllPort(iplist)
	// log.Println(allPortOk)
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
