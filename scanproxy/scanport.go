package scanproxy

import (
	"fmt"
	"log"
	"net"
	"time"

	tcp "github.com/tevino/tcp-shaker"
)

var (
	checkPortTimeout = 2 * time.Second
	portMax          = 65535
	step             = 25
)

//checkPort 查询端口是否开放
func checkPort(ipstr string, port int, ch chan map[string]int) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		log.Println("error IP format:", ipstr)
		ch <- nil
	}
	// tcpAddr := net.TCPAddr{
	// 	IP:   ip,
	// 	Port: port,
	// }
	address := fmt.Sprintf("%v:%v", ip, port)
	conn, err := net.DialTimeout("tcp", address, checkPortTimeout) //("tcp", nil, &tcpAddr)
	if err == nil {
		log.Println("open IP & port", ip, port)
		conn.Close()
		ch <- map[string]int{string(ipstr): port}
	} else {
		// log.Println("scan port fail：", ipstr, port)
		if conn != nil {
			conn.Close()
		}
		ch <- nil
	}
}

//checkPortBySyn syn方式查询端口是否开放，仅支持linux2.4+
func checkPortBySyn(ipstr string, port int, ch chan map[string]int) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		log.Println("error IP format:", ipstr)
		ch <- nil
	}
	address := fmt.Sprintf("%v:%v", ip, port)
	checker := tcp.NewChecker(true)
	if err := checker.InitChecker(); err != nil {
		log.Fatal("Checker init failed:", err)
		ch <- nil
	}
	err := checker.CheckAddr(address, checkPortTimeout)
	switch err {
	case tcp.ErrTimeout:
		// fmt.Println("Connect to host timed out")
		// log.Println(ipstr, port, err)
		ch <- nil
	case nil:
		log.Println("open IP & port", ip, port)
		ch <- map[string]int{ipstr: port}
	default:
		// if e, ok := err.(*tcp.ErrConnect); ok {
		// 	fmt.Println("Connect to host failed:", e)
		// } else {
		// 	fmt.Println("Error occurred while connecting:", err)
		// }
		ch <- nil
	}
}

func scanPort(iplist *[]string, startPort int, stepMax int) (*[]map[string]int, int) {
	ch := make(chan map[string]int, 1000)
	var portOkList []map[string]int
	var value map[string]int
	i := 0
	//分阶段扫描端口
	for n := startPort; n <= stepMax; n++ {
		//循环处理IP段
		log.Println("scan port:", n)
		for j := len(*iplist) - 1; j > 0; j-- {
			// log.Println(iplist[j], i, ch)
			go checkPortBySyn((*iplist)[j], n, ch)
			time.Sleep(1 * time.Millisecond)
			i++
		}
	}
	//分阶段回收被BLOCK的协程
	// step := stepMax - startPort
	// lenList := (len(*iplist) - 2) * step
	time.Sleep(2 * time.Second)
	for m := 1; m <= i; m++ {
		// for value := range ch {
		value = <-ch
		// log.Println(value)
		if value != nil {
			portOkList = append(portOkList, value)
		}
	}
	close(ch)
	return &portOkList, stepMax
}

//ScanAllPort 分段扫描65535全部端口
func ScanAllPort(iplist *[]string) *[]map[string]int {
	var stepMax int
	var portOkList []map[string]int
	for i := 1; i <= portMax; i++ {
		if (step + i) > portMax {
			stepMax = portMax
		} else {
			stepMax = step + i
		}
		portOpen, endPort := scanPort(iplist, i, stepMax)
		i = endPort
		if portOpen != nil {
			portOkList = append(portOkList, (*portOpen)...)
		}
		// log.Println(portOkList)
	}
	return &portOkList
}

//InternetAllScan 全部IP或指定区域IP扫描全端口
func InternetAllScan(area string) {
	totalPage := 1
	var ipmap []map[string]string
	var err error
	for i := 1; i <= totalPage; i++ {
		ipmap, _, totalPage, err = GetApnicIP(area, i, 1)
		if err != nil {
			log.Fatalln(err)
		}
		startip := ipmap[0]["startip"]
		getArea := ipmap[0]["area"]
		log.Println("start scan IP:", startip)
		iplist := formatInternetIPList(startip)
		portOpenList := ScanAllPort(iplist)
		go func() {
			httpProxy := CheckHTTPForList(portOpenList)
			socksProxy := CheckSocksForList(portOpenList)
			allproxyList := append((*httpProxy), (*socksProxy)...)
			if allproxyList != nil {
				SaveProxy(&allproxyList, getArea)
			}
		}()
	}
}

//InternetFastScan 常用代理端口快速扫描
// func InternetFastScan(area string) {

// }
