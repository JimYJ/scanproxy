package scanproxy

import (
	"errors"
	"log"
	"time"

	"github.com/JimYJ/scanproxy/config"
)

func saveProxy(proxyList *[]map[string]string, area interface{}) (bool, error) {
	if proxyList == nil {
		return false, errors.New("proxyList is nil")
	}
	mysqlDB := config.MySQL()
	rollBack := false
	var err2 error
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	var proxyExist = make([]int64, len(*proxyList))
	for i := 0; i < len(*proxyList); i++ {
		checkExist, err := mysqlDB.GetVal("select id from proxyip where ip = ? and port = ? and protocol = ?", (*proxyList)[i]["ip"], (*proxyList)[i]["port"], (*proxyList)[i]["protocol"])
		if err != nil || checkExist == nil {
			proxyExist[i] = 0
		} else {
			id, ok := checkExist.(int64)
			if ok {
				proxyExist[i] = id
			} else {
				proxyExist[i] = 0
			}
		}
	}
	tx, _ := mysqlDB.Begin()
	for i := 0; i < len(*proxyList); i++ {
		if proxyExist[i] == 0 {
			_, err2 = tx.Insert("insert into proxyip set ip = ?,port = ?,protocol = ?,area = ?,status = ?,createtime = ?,updatetime = ?", (*proxyList)[i]["ip"], (*proxyList)[i]["port"], (*proxyList)[i]["protocol"], area, 1, nowTime, nowTime)
		} else {
			_, err2 = tx.Update("update proxyip set updatetime = ? where id = ?", nowTime, proxyExist[i])
		}
		if err2 != nil {
			rollBack = true
		}
	}
	if rollBack {
		tx.Rollback()
		return false, err2
	}
	tx.Commit()
	log.Println("save proxy:", proxyList)
	return true, nil
}
