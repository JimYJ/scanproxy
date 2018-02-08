package scanproxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	timeouts    = 5
	testWeb     = "https://email.163.com"
	testKeyWord = "网易免费邮箱"
)

//CheckHTTP 测试是否是HTTP代理服务器
func CheckHTTP(ip string, port int, protocol string) bool {
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
		defer resp.Body.Close()
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil && strings.Contains(string(body), testKeyWord) {
					return true
				}
			}
		}
	}
	return false
}
