package scanproxy

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/JimYJ/go-queue"

	tcp "github.com/tevino/tcp-shaker"
)

var (
	checkPortTimeout   = 5 * time.Second
	portMax            = 65535
	step               = 1
	once               sync.Once
	queueMaxConcurrent = 1000
	fastScanPort       = [...]int{3128, 8000, 8888, 8080, 8088, 1080, 9000, 80, 8118, 53281, 54566, 808, 443, 8081, 8118, 65103, 3333, 45619, 65205, 45619, 55379, 65535, 2855, 10200, 22722, 64334, 3654, 53124, 5433}
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

//checkPortBySynForQueue syn方式查询端口是否开放，仅支持linux2.4+ 队列专用函数
func checkPortBySynForQueue(value ...interface{}) {
	if len(value) < 3 {
		return
	}
	ipstr := value[0].(string)
	port := value[1].(int)
	ch := value[2].(chan map[string]int)
	ip := net.ParseIP(ipstr)
	if ip == nil {
		log.Println("error IP format:", ipstr)
		ch <- nil
	}
	address := fmt.Sprintf("%v:%v", ip, port)
	checker := tcp.NewChecker(true)
	if err := checker.InitChecker(); err != nil {
		// log.Fatal("Checker init failed:", err)
		ch <- nil
	}
	err := checker.CheckAddr(address, checkPortTimeout)
	switch err {
	case tcp.ErrTimeout:
		// fmt.Println("Connect to host timed out")
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
	ch := make(chan map[string]int, 2000)
	var portOkList []map[string]int
	var value map[string]int
	i := 0
	//分阶段扫描端口
	for n := startPort; n <= stepMax; n++ {
		//循环处理IP段
		// log.Println("scan port:", n)
		for j := len(*iplist) - 1; j >= 0; j-- {
			go checkPortBySyn((*iplist)[j], n, ch)
			time.Sleep(1 * time.Millisecond)
			i++
		}
	}
	//分阶段回收被BLOCK的协程
	// log.Println(i)
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

//scanAllPort 分段扫描65535全部端口
func scanAllPort(iplist *[]string) *[]map[string]int {
	var stepMax int
	var portOkList []map[string]int
	for i := 0; i < len(fastScanPort); i++ {
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

//scanFastPort 快速扫描常用端口,将请求发送到队列
func scanFastPort(iplist *[]string, getArea string, ch chan map[string]int) {
	listenQueueResults(ch, getArea)
	for n := 0; n <= len(fastScanPort)-1; n++ {
		// log.Println("scan port:", n)
		for j := len(*iplist) - 1; j >= 0; j-- {
			job := new(queue.Job)
			job.ID = time.Now().Local().UnixNano()
			job.FuncQueue = checkPortBySynForQueue
			job.Payload = []interface{}{(*iplist)[j], fastScanPort[n], ch}
			queue.JobQueue <- job
			// time.Sleep(1 * time.Second)
		}
	}
}

func listenQueueResults(ch chan map[string]int, getArea string) {
	once.Do(func() {
		go func() {
			for {
				// for value := range ch {
				var portOkList []map[string]int
				for i := 0; i < queueMaxConcurrent; i++ {
					value := <-ch
					// log.Println(value)
					if value != nil {
						if value != nil {
							portOkList = append(portOkList, value)
						}
					}
				}
				if portOkList != nil {
					go func() {
						httpProxy := checkHTTPForList(&portOkList)
						socksProxy := checkSocksForList(&portOkList)
						allproxyList := append((*httpProxy), (*socksProxy)...)
						if allproxyList != nil {
							saveProxy(&allproxyList, getArea)
						}
					}()
				}
			}
		}()
	})
}

//InternetAllScan 全部IP或指定区域IP扫描全端口
func InternetAllScan(area string, ipStep int) {
	totalPage := 1
	var ipmap *[]map[string]string
	var err error
	var iplist []string
	for i := 1; i <= totalPage; i++ {
		iplist = make([]string, 0)
		ipmap, _, totalPage, err = getApnicIP(area, i, ipStep)
		if err != nil {
			log.Fatalln(err)
		}
		getArea := (*ipmap)[0]["area"]
		for m := 0; m < len(*ipmap); m++ {
			startip := (*ipmap)[m]["startip"]
			log.Println("start scan IP:", startip)
			iplist = append(iplist, formatInternetIPList(startip)...)
		}
		portOpenList := scanAllPort(&iplist)
		go func() {
			httpProxy := checkHTTPForList(portOpenList)
			socksProxy := checkSocksForList(portOpenList)
			allproxyList := append((*httpProxy), (*socksProxy)...)
			if allproxyList != nil {
				saveProxy(&allproxyList, getArea)
			}
		}()
	}
}

//InternetFastScan 常用代理端口快速扫描
func InternetFastScan(area string, ipStep int) {
	totalPage := 1
	var ipmap *[]map[string]string
	var err error
	var iplist []string
	queue.InitQueue(queueMaxConcurrent, false)
	ch := make(chan map[string]int, queueMaxConcurrent)
	for i := 1; i <= totalPage; i++ {
		iplist = nil
		for j := 0; j < 10; j++ {
			ipmap, _, totalPage, err = getApnicIP(area, i, ipStep)
			if err != nil {
				log.Println(err)
			} else {
				break
			}
			time.Sleep(10 * time.Second)
		}
		getArea := (*ipmap)[0]["area"]
		for m := 0; m < len(*ipmap); m++ {
			startip := (*ipmap)[m]["startip"]
			log.Println("start fast scan IP:", startip)
			iplist = append(iplist, formatInternetIPList(startip)...)
		}
		scanFastPort(&iplist, getArea, ch)
		time.Sleep(15 * time.Second)
	}
}

//SetQueueMaxConcurrent 设置队列并发数
func SetQueueMaxConcurrent(maxConcurrent string) {
	mc, err := strconv.Atoi(maxConcurrent)
	if err != nil || mc < 1 {
		return
	}
	queueMaxConcurrent = mc
}
