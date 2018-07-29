package main

import (
	"flag"
	"log"
	"scanproxy/scanproxy"

	"github.com/JimYJ/easysql/mysql"
)

var (
	dbhost        = ""
	dbport        = 3306
	dbname        = ""
	dbuser        = ""
	dbpass        = ""
	charset       = "utf8mb4"
	maxIdleConns  = 500
	maxOpenConns  = 500
	ipStep        = 10
	area          string
	ipstep        int
	mode          bool
	maxConcurrent int
)

func main() {
	initDBConn()
	if mode {
		log.Println("now work in fast mode,ip scan step:", ipstep, "area:", area, "maxConcurrent:", maxConcurrent)
		scanproxy.SetQueueMaxConcurrent(maxConcurrent)
		scanproxy.InternetFastScan(area, ipstep)
	} else {
		log.Println("now work in normal mode,ip scan step:", ipstep, "area:", area)
		scanproxy.InternetAllScan(area, ipstep)
	}
}

func init() {
	flag.BoolVar(&mode, "f", false, "scan mode, fast is only scan common port, default is scan all port(shorthand)")
	flag.IntVar(&ipstep, "i", 10, "set scan how many IP segment in same times, it will affect memory footprint(shorthand)")
	flag.IntVar(&maxConcurrent, "m", 200, "maximum concurrency number(shorthand)")
	flag.StringVar(&area, "a", "CN", "country codes, see ISO 3166-1(shorthand)")
	flag.Parse()
}

func initDBConn() {
	mysql.Init(dbhost, dbport, dbname, dbuser, dbpass, charset, maxIdleConns, maxOpenConns)
	_, err := mysql.GetMysqlConn()
	if err != nil {
		log.Panic(err)
	}
}
