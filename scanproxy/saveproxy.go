package scanproxy

import (
	"errors"
	"github.com/JimYJ/easysql/mysql"
	"log"
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
	for i := 0; i < len(*proxyList); i++ {
		_, err2 = mysqlDB.TxInsert(mysql.Statement, "insert into proxyip set ip = ?,port = ?,protocol = ?,area = ?", (*proxyList)[i]["ip"], (*proxyList)[i]["port"], (*proxyList)[i]["protocol"], area)
		if err2 != nil {
			rollBack = true
		}
	}
	if rollBack {
		mysqlDB.TxRollback()
		return false, err2
	}
	mysqlDB.TxCommit()
	return true, nil
}
