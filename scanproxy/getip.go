package scanproxy

import (
	"errors"
	"log"
	"strconv"

	"github.com/JimYJ/easysql/mysql"
)

var (
	//ipCount 可扫描IP总数
	ipCount = 0
)

//GetApnicIP 获取可扫描IP列表
func GetApnicIP(area string, curPage int, prePage int) ([]map[string]string, int, int, error) {
	mysqlConn, err := mysql.GetMysqlConn()
	if err != nil {
		return nil, 0, 0, err
	}
	paginate, total, totalPage := paginate(area, curPage, prePage)
	if total == 0 || totalPage == 0 {
		return nil, 0, 0, errors.New("area is error")
	}
	var query string
	if area != "" {
		query = "select startip,endip from apniciplib where area = ?" + paginate
	} else {
		query = "select startip,endip from apniciplib" + paginate
	}
	iplist, err := mysqlConn.GetResults(mysql.Statement, query, area)
	if err != nil {
		return nil, 0, 0, err
	}
	return iplist, total, totalPage, nil
}

func paginate(area string, curPage int, prePage int) (string, int, int) {
	if ipCount == 0 {
		GetIPCount(area)
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
		return "limit 0," + string(prePage), ipCount, totalPage
	}
	start := (curPage - 1) * prePage
	return "limit " + string(start) + "," + string(prePage), ipCount, totalPage
}

func getTotalPage(prePage int) int {
	totalPage := ipCount / prePage
	if ipCount%prePage != 0 {
		totalPage++
	}
	return totalPage
}

//GetIPCount 获取IP总数
func GetIPCount(area string) error {
	mysqlConn, err := mysql.GetMysqlConn()
	if err != nil {
		return err
	}
	var query string
	if area != "" {
		query = "select count(id) as count from apniciplib where area = ?"
	} else {
		query = "select count(id) as count from apniciplib"
	}
	count, err := mysqlConn.GetVal(mysql.Statement, query, area)
	if err != nil {
		return err
	}
	ipCount, _ = strconv.Atoi(count)
	return nil
}

//GetIPtemp 获取内网IP列表
func GetIPtemp() []string {
	var a int
	var iplist = make([]string, 255)
	for i := 1; i < 256; i++ {
		a = i - 1
		iplist[a] = "192.168.10." + strconv.Itoa(i)
	}
	return iplist
}
