package scanproxy

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/JimYJ/easysql/mysql"
)

func saveProxy(proxyList *[]map[string]string, area string) (bool, error) {
	if proxyList == nil {
		return false, errors.New("proxyList is nil")
	}
	mysqlDB, err := mysql.GetMysqlConn()
	if err != nil {
		log.Panic(err)
	}
	rollBack := false
	var err2 error
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	var proxyExist = make([]int, len(*proxyList))
	for i := 0; i < len(*proxyList); i++ {
		checkExist, err := mysqlDB.GetVal(mysql.Statement, "select id from proxyip where ip = ? and port = ? and protocol = ?", (*proxyList)[i]["ip"], (*proxyList)[i]["port"], (*proxyList)[i]["protocol"])
		if err != nil || checkExist == "" {
			proxyExist[i] = 0
		} else {
			id, err := strconv.Atoi(checkExist)
			if err == nil {
				proxyExist[i] = id
			} else {
				proxyExist[i] = 0
			}
		}
	}
	mysqlDB.TxBegin()
	for i := 0; i < len(*proxyList); i++ {
		if proxyExist[i] == 0 {
			_, err2 = mysqlDB.TxInsert(mysql.Statement, "insert into proxyip set ip = ?,port = ?,protocol = ?,area = ?,status = ?,createtime = ?,updatetime = ?", (*proxyList)[i]["ip"], (*proxyList)[i]["port"], (*proxyList)[i]["protocol"], area, 1, nowTime, nowTime)
		} else {
			_, err2 = mysqlDB.TxUpdate(mysql.Statement, "update proxyip set updatetime = ? where id = ?", nowTime, proxyExist[i])
		}
		if err2 != nil {
			rollBack = true
		}
	}
	if rollBack {
		mysqlDB.TxRollback()
		return false, err2
	}
	mysqlDB.TxCommit()
	log.Println("save proxy:", proxyList)
	return true, nil
}
