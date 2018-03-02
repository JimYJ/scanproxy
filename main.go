package main

import (
	"log"
	"os"
	"scanproxy/scanproxy"
	"strconv"

	"github.com/JimYJ/easysql/mysql"
)

var (
	dbhost       = "rm-bp18iy73784671903yo.mysql.rds.aliyuncs.com"
	dbport       = 3306
	dbname       = "dutyfree"
	dbuser       = "root_xw"
	dbpass       = "Xw_19920602_wX"
	charset      = "utf8mb4"
	maxIdleConns = 500
	maxOpenConns = 500
	ipStep       = 10
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
	// iplist := make([]map[string]int, 1)
	// a := make(map[string]int)
	// a["133.130.103.208"] = 8080
	// iplist = append(iplist, a)
	// b := scanproxy.CheckHTTPForList(&iplist)
	// log.Println(b)
	// log.Println(scanproxy.CheckSocksForList(&iplist))
	// log.Println(scanproxy.CheckSocks5("23.254.153.205", 25357, "tcp"))
	// c, err := scanproxy.SaveProxy(b, "JP")
	// log.Println(c, err)
	mode = os.Args[1]
	ipstep, err = strconv.Atoi(os.Args[2])
	if err != nil {
		ipstep = IPSipSteptep
	}

	initDBConn()
	scanproxy.InternetAllScan("CN", ipstep)

}

func initDBConn() {
	mysql.Init(dbhost, dbport, dbname, dbuser, dbpass, charset, maxIdleConns, maxOpenConns)
	_, err := mysql.GetMysqlConn()
	if err != nil {
		log.Panic(err)
	}
}
