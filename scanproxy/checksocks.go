package scanproxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
	"h12.me/socks"
)

//CheckSocks 测试是否是SOCKS代理，支持SOCKS4,SOCKS4a,SOCKS5
func CheckSocks(ip string, port int, protocol int) bool {
	strURL := fmt.Sprintf("%v:%v", ip, port)
	client := http.Client{
		Transport: &http.Transport{
			Dial: socks.DialSocksProxy(protocol, strURL),
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
		log.Println(err)
		if resp != nil {
			defer resp.Body.Close()
		}
	}
	return false
}

//CheckSocks5 测试是否是SOCKS5代理
func CheckSocks5(ip string, port int, protocol string) bool {
	strURL := fmt.Sprintf("%v:%v", ip, port)
	dialer, err := proxy.SOCKS5(protocol, strURL, nil, proxy.Direct)
	if err == nil {
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial
		resp, err := httpClient.Get(testWeb)
		if resp != nil {
			defer resp.Body.Close()
		}
		log.Println(err, resp)
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

func checkSocksForList(iplist *[]map[string]int) *[]map[string]string {
	var proxyOK []map[string]string
	for i := 0; i < len(*iplist); i++ {
		for k, v := range (*iplist)[i] {
			if CheckSocks(k, v, socks.SOCKS4) {
				proxy := map[string]string{"ip": k, "port": strconv.Itoa(v), "protocol": "socks4"}
				proxyOK = append(proxyOK, proxy)
			}
			if CheckSocks(k, v, socks.SOCKS4A) {
				proxy := map[string]string{"ip": k, "port": strconv.Itoa(v), "protocol": "socks4a"}
				proxyOK = append(proxyOK, proxy)
			}
			if CheckSocks(k, v, socks.SOCKS5) {
				proxy := map[string]string{"ip": k, "port": strconv.Itoa(v), "protocol": "socks5"}
				proxyOK = append(proxyOK, proxy)
			}
		}
	}
	return &proxyOK
}
