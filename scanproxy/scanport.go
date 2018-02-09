package scanproxy

import (
	"fmt"
	"log"
	"net"
	"time"
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
	conn, err := net.DialTimeout("tcp", address, 30*time.Second) //("tcp", nil, &tcpAddr)
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
