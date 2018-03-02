package scanproxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	timeouts    = 10
	testWeb     = "http://auto.163.com"
	testKeyWord = "auto.163.com"
)

//checkHTTP 测试是否是HTTP代理服务器
func checkHTTP(ip string, port int, protocol string) bool {
	strURL := fmt.Sprintf("%v://%v:%v", protocol, ip, port)
	proxyURL, err := url.Parse(strURL)
	if err == nil {
		client := http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: time.Duration(timeouts) * time.Second,
		}
		resp, err := client.Get(testWeb)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil && strings.Contains(string(body), testKeyWord) {
					return true
				}
			}
		} else {
			if resp != nil {
				defer resp.Body.Close()
			}
			log.Println(err)
		}
	}
	return false
}

func checkHTTPForList(iplist *[]map[string]int) *[]map[string]string {
	var proxyOK []map[string]string
	for i := 0; i < len(*iplist); i++ {
		for k, v := range (*iplist)[i] {
			if checkHTTP(k, v, "http") {
				proxy := map[string]string{"ip": k, "port": strconv.Itoa(v), "protocol": "http"}
				proxyOK = append(proxyOK, proxy)
			}
			if checkHTTP(k, v, "https") {
				proxy := map[string]string{"ip": k, "port": strconv.Itoa(v), "protocol": "https"}
				proxyOK = append(proxyOK, proxy)
			}
		}
	}
	return &proxyOK
}
