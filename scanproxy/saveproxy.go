package scanproxy

import (
	"errors"
	"log"
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
	mysqlDB.TxBegin()
	rollBack := false
	var err2 error
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	for i := 0; i < len(*proxyList); i++ {
		_, err2 = mysqlDB.TxInsert(mysql.Statement, "insert into proxyip set ip = ?,port = ?,protocol = ?,area = ?,status = ?,createtime = ?,updatetime = ?", (*proxyList)[i]["ip"], (*proxyList)[i]["port"], (*proxyList)[i]["protocol"], area, 1, nowTime, nowTime)
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
