package scanproxy

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/JimYJ/easysql/mysql"
)

var (
	//ipCount 可扫描IP总数
	ipCount = 0
)

//getApnicIP 获取可扫描IP列表
func getApnicIP(area string, curPage int, prePage int) (*[]map[string]string, int, int, error) {
	mysqlConn, err := mysql.GetMysqlConn()
	if err != nil {
		return nil, 0, 0, err
	}
	paginate, total, totalPage := paginate(area, curPage, prePage)
	if total == 0 || totalPage == 0 {
		return nil, 0, 0, errors.New("area is error")
	}
	var query string
restart:
	recordID, IPID := getRecord(mysqlConn, area)
	if area != "" {
		query = "select id,startip,area from apniciplib where area = ? and id > ?" + paginate
	} else {
		query = "select id,startip,area from apniciplib id < ?" + paginate
	}
	iplist, err := mysqlConn.GetResults(mysql.Statement, query, area, IPID)
	if err != nil {
		return nil, 0, 0, err
	}
	if len(iplist) == 0 {
		saveRecord(mysqlConn, 0, recordID, "", area)
		goto restart
	}
	idstr := iplist[len(iplist)-1]["id"]
	startIP := iplist[len(iplist)-1]["startip"]
	id, err2 := strconv.Atoi(idstr)
	if err2 == nil {
		saveRecord(mysqlConn, id, recordID, startIP, area)
	} else {
		log.Println(err2)
	}
	return &iplist, total, totalPage, nil
}

func paginate(area string, curPage int, prePage int) (string, int, int) {
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
		return " limit 0," + strconv.Itoa(prePage), ipCount, totalPage
	}
	start := (curPage - 1) * prePage
	return " limit " + strconv.Itoa(start) + "," + strconv.Itoa(prePage), ipCount, totalPage
}

func getTotalPage(prePage int) int {
	totalPage := ipCount / prePage
	if ipCount%prePage != 0 {
		totalPage++
	}
	return totalPage
}

//getIPCount 获取IP总数
func getIPCount(area string) error {
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

//getIPLocalNetwork 获取内网IP列表
func getIPLocalNetwork() []string {
	var a int
	var iplist = make([]string, 255)
	for i := 1; i < 256; i++ {
		a = i - 1
		iplist[a] = "192.168.10." + strconv.Itoa(i)
	}
	return iplist
}

func formatInternetIPList(ipsatrt string) []string {
	var iplist = make([]string, 0)
	b := strings.Split(ipsatrt, ".")
	c := strings.Join(b[0:len(b)-1], ".")
	for i := 1; i <= 254; i++ {
		iplist = append(iplist, c+"."+strconv.Itoa(i))
	}
	return iplist
}

func saveRecord(mysqlConn *mysql.MysqlDB, id int, recordID int, startip string, area string) {
	var query string
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	var err error
	for i := 0; i < 3; i++ {
		if recordID == 0 {
			query = "insert scanrecord set ipid = ?,area = ?,createtime = ?,updatetime = ?,startip = ?"
			_, err = mysqlConn.Insert(mysql.Statement, query, id, area, nowTime, nowTime, startip)
		} else {
			query = "update scanrecord set ipid = ?,updatetime = ?,startip = ? where id =?"
			_, err = mysqlConn.Update(mysql.Statement, query, id, nowTime, startip, recordID)
		}
		if err != nil {
			log.Println(err)
		} else {
			break
		}
	}
}

func getRecord(mysqlConn *mysql.MysqlDB, area string) (int, int) {
	query := "select id,ipid from scanrecord where area = ?"
	var rs map[string]string
	var err error
	for i := 0; i < 3; i++ {
		err = nil
		rs, err = mysqlConn.GetRow(mysql.Statement, query, area)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		return 0, 0
	}
	idstr := rs["id"]
	idipstr := rs["ipid"]
	id, err2 := strconv.Atoi(idstr)
	ipid, err3 := strconv.Atoi(idipstr)
	if err2 != nil || err3 != nil {
		return 0, 0
	}
	return id, ipid
}
