package main

import (
	"log"
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
	// runtime.GOMAXPROCS(runtime.NumCPU())
	// initDBConn()
	// ch := make(chan map[string]int, 1)
	// go scanproxy.CheckPort("192.168.10.242", 80, ch)
	// log.Println(<-ch)
	iplist := scanproxy.GetIPtemp()
	allPortOk := scanAllPort(iplist)
	log.Println(allPortOk)
}

func initDBConn() {
	mysql.Init(dbhost, dbport, dbname, dbuser, dbpass, charset, maxIdleConns, maxOpenConns)
	_, err := mysql.GetMysqlConn()
	if err != nil {
		log.Panic(err)
	}
}

func scanPort(iplist []string, startPort int, stepMax int) ([]map[string]int, int) {
	ch := make(chan map[string]int, 1000)
	var portOkList []map[string]int
	var value map[string]int
	//分阶段扫描端口
	for n := startPort; n <= stepMax; n++ {
		//循环处理IP段
		// log.Println("scan port:", i)
		for j := len(iplist) - 1; j > 0; j-- {
			// log.Println(iplist[j], i, ch)
			go scanproxy.CheckPort(iplist[j], n, ch)
			time.Sleep(1 * time.Millisecond)

		}
	}
	//分阶段回收被BLOCK的协程
	step := stepMax - startPort
	for m := 0; m <= ((len(iplist) - 2) * step); m++ {
		// for value := range ch {
		value = <-ch
		// log.Println(value)
		if value != nil {
			portOkList = append(portOkList, value)
		}
	}
	time.Sleep(1 * time.Second)
	close(ch)
	return portOkList, stepMax
}

func scanAllPort(iplist []string) []map[string]int {
	step := 25
	var stepMax, endPort int
	var portOkList []map[string]int
	for i := 1; i <= portMax; i++ {
		if (step + i) > portMax {
			stepMax = portMax
		} else {
			stepMax = step + i
		}
		portOkList, endPort = scanPort(iplist, i, stepMax)
		i = endPort
		// log.Println(portOkList)
	}
	return portOkList
}
