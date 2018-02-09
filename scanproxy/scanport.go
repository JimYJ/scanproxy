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
)

//CheckPort 查询端口是否开放
func CheckPort(ipstr string, port int, ch chan map[string]int) {
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

//CheckPortBySyn syn方式查询端口是否开放，仅支持linux2.4+
func CheckPortBySyn(ipstr string, port int, ch chan map[string]int) {
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
		ch <- nil
	case nil:
		log.Println("open IP & port", ip, port)
		ch <- map[string]int{string(ipstr): port}
	default:
		// if e, ok := err.(*tcp.ErrConnect); ok {
		// 	fmt.Println("Connect to host failed:", e)
		// } else {
		// 	fmt.Println("Error occurred while connecting:", err)
		// }
		ch <- nil
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
			go CheckPortBySyn(iplist[j], n, ch)
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

//ScanAllPort 分段扫描65535全部端口
func ScanAllPort(iplist []string) []map[string]int {
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
