package main

import (
	"github.com/JimYJ/scanproxy/scanproxy"
	"log"
	"os"
	"strconv"

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
	ipStep       = 10
	area         string
)

func main() {
	var ipstep int
	var err error
	var mode string
	var maxConcurrent string
	if len(os.Args) >= 2 {
		mode = os.Args[1]
	} else {
		mode = ""
	}
	if len(os.Args) >= 3 {
		ipstep, err = strconv.Atoi(os.Args[2])
		if err != nil {
			ipstep = ipStep
		}
	} else {
		ipstep = ipStep
	}
	if len(os.Args) >= 4 {
		area = os.Args[3]
	} else {
		area = ""
	}
	initDBConn()
	if mode == "-f" {
		if len(os.Args) >= 5 {
			maxConcurrent = os.Args[4]
		} else {
			maxConcurrent = "200"
		}
		log.Println("now work in fast mode,ip scan step:", ipstep, "area:", area, "maxConcurrent:", maxConcurrent)
		scanproxy.SetQueueMaxConcurrent(maxConcurrent)
		scanproxy.InternetFastScan(area, ipstep)
	} else {
		log.Println("now work in normal mode,ip scan step:", ipstep, "area:", area)
		scanproxy.InternetAllScan(area, ipstep)
	}
}

func initDBConn() {
	mysql.Init(dbhost, dbport, dbname, dbuser, dbpass, charset, maxIdleConns, maxOpenConns)
	_, err := mysql.GetMysqlConn()
	if err != nil {
		log.Panic(err)
	}
}
