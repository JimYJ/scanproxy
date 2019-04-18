package scanproxy

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/JimYJ/scanproxy/config"
)

var (
	//ipCount 可扫描IP总数
	ipCount int64 = 0
)

//getApnicIP 获取可扫描IP列表
func getApnicIP(area string, curPage int, prePage int) (*[]map[string]interface{}, int64, int, error) {
	mysqlDB := config.MySQL()
	paginate, total, totalPage := paginate(area, curPage, prePage)
	if total == 0 || totalPage == 0 {
		return nil, 0, 0, errors.New("area is error")
	}
	var query string
restart:
	recordID, IPID := getRecord(area)
	if area != "" {
		query = fmt.Sprintf("select id,startip,area from apniciplib where area = ? and id > ? %s", paginate)
	} else {
		query = fmt.Sprintf("select id,startip,area from apniciplib where id < ? %s", paginate)
	}
	iplist, err := mysqlDB.GetResults(query, area, IPID)
	if err != nil {
		return nil, 0, 0, err
	}
	if len(iplist) == 0 {
		saveRecord(0, recordID, "", area)
		goto restart
	}
	startIP := iplist[len(iplist)-1]["startip"]
	id, ok := iplist[len(iplist)-1]["id"].(int64)
	if !ok {
		saveRecord(id, recordID, startIP, area)
	} else {
		log.Println("get ip list error")
	}
	return &iplist, total, totalPage, nil
}

func paginate(area string, curPage int, prePage int) (string, int64, int) {
	if ipCount == 0 {
		getIPCount(area)
		if ipCount == 0 {
			log.Println("area is error!")
			return "", 0, 0
		}
	}
	totalPage := getTotalPage(prePage)
	if curPage > totalPage {
		curPage = totalPage
	}
	if curPage == 0 || curPage == 1 {
		return fmt.Sprintf("limit 0,%d", prePage), ipCount, totalPage
	}
	start := (curPage - 1) * prePage
	return "limit " + strconv.Itoa(start) + "," + strconv.Itoa(prePage), ipCount, totalPage
}

func getTotalPage(prePage int) int {
	totalPage := int(ipCount) / prePage
	if int(ipCount)%prePage != 0 {
		totalPage++
	}
	return totalPage
}

//getIPCount 获取IP总数
func getIPCount(area string) error {
	mysqlConn := config.MySQL()
	var query string
	if area != "" {
		query = "select count(id) as count from apniciplib where area = ?"
	} else {
		query = "select count(id) as count from apniciplib"
	}
	count, err := mysqlConn.GetVal(query, area)
	if err != nil {
		return err
	}
	ipCount = count.(int64)
	return nil
}

//getIPLocalNetwork 获取内网IP列表
func getIPLocalNetwork() []string {
	var a int
	var iplist = make([]string, 255)
	for i := 1; i < 256; i++ {
		a = i - 1
		iplist[a] = fmt.Sprintf("192.168.10.%d", i)
	}
	return iplist
}

func formatInternetIPList(ips interface{}) []string {
	ipsatrt, ok := ips.(string)
	if !ok {
		return nil
	}
	var iplist = make([]string, 0)
	b := strings.Split(ipsatrt, ".")
	c := strings.Join(b[0:len(b)-1], ".")
	for i := 1; i <= 254; i++ {
		iplist = append(iplist, c+"."+strconv.Itoa(i))
	}
	return iplist
}

func saveRecord(id, recordID int64, startip interface{}, area string) {
	mysqlConn := config.MySQL()
	var query string
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	var err error
	for i := 0; i < 3; i++ {
		if recordID == 0 {
			query = "insert scanrecord set ipid = ?,area = ?,createtime = ?,updatetime = ?,startip = ?"
			_, err = mysqlConn.Insert(query, id, area, nowTime, nowTime, startip)
		} else {
			query = "update scanrecord set ipid = ?,updatetime = ?,startip = ? where id =?"
			_, err = mysqlConn.Update(query, id, nowTime, startip, recordID)
		}
		if err != nil {
			log.Println(err)
		} else {
			break
		}
	}
}

func getRecord(area string) (int64, int64) {
	mysqlConn := config.MySQL()
	query := "select id,ipid from scanrecord where area = ?"
	var rs map[string]interface{}
	var err error
	for i := 0; i < 3; i++ {
		err = nil
		rs, err = mysqlConn.GetRow(query, area)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		return 0, 0
	}
	id, ok := rs["id"].(int64)
	ipid, ok2 := rs["ipid"].(int64)
	if !ok || !ok2 {
		return 0, 0
	}
	return id, ipid
}
