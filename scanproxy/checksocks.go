package scanproxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/proxy"
	"h12.me/socks"
)

//CheckSocks 测试是否是SOCKS代理，支持SOCKS4，SOCKS4a,SOCKS5
func CheckSocks(ip string, port int, protocol int) bool {
	strURL := fmt.Sprintf("%v://%v:%v", protocol, ip, port)
	client := http.Client{
		Transport: &http.Transport{
			Dial: socks.DialSocksProxy(protocol, strURL),
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
