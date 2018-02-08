package main

import (
	"log"
	"runtime"
	"scanproxy/scanproxy"
	"time"

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
	portMax      = 65535
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// initDBConn()
	// ch := make(chan map[string]int, 1)
	// go scanproxy.CheckPort("192.168.10.242", 80, ch)
	// log.Println(<-ch)
	iplist := scanproxy.GetIPtemp()
	scanAllPort(iplist)
}

func initDBConn() {
	mysql.Init(dbhost, dbport, dbname, dbuser, dbpass, charset, maxIdleConns, maxOpenConns)
	_, err := mysql.GetMysqlConn()
	if err != nil {
		log.Panic(err)
	}
}

func scanAllPort(iplist []string) {
	ch := make(chan map[string]int, 3000)
	var portOkList []map[string]int
	var stepmax int
	max := 65535
	step := 25
	var results map[string]int
	for i := 1; i <= max; i++ {
		if (i + step) < max {
			stepmax = i + step
		} else {
			stepmax = max
		}
		//分阶段扫描端口
		for n := i; n <= stepmax; n++ {
			//循环处理IP段
			for j := len(iplist) - 1; j > 0; j-- {
				// log.Println(iplist[j], i, ch)
				go scanproxy.CheckPort(iplist[j], i, ch)
			}
			i = n
		}
		//分阶段回收被BLOCK的协程
		for m := 0; m <= ((len(iplist) - 2) * step); m++ {
			results = <-ch
			if results != nil {
				portOkList = append(portOkList, results)
			}
		}
		time.Sleep(1 * time.Second)
	}
	log.Println(portOkList)
}
