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
