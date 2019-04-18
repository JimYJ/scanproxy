package config

import (
	"io/ioutil"
	"log"
	"sync"

	mysql "github.com/JimYJ/easysql/mysql/v2"
	cache "github.com/patrickmn/go-cache"

	"gopkg.in/yaml.v2"
)

var (
	once sync.Once
	c    *cache.Cache
)

// 初始化参数
var (
	Host, User, Pass, Name string
	Port                   int
)

// Config 基础配置
type config struct {
	MySQL Mysql
}

// Mysql 数据库配置
type Mysql struct {
	Host, User, Pass, Name string
	Port                   int
}

func init() {
	configInit()
	MySQL()
}

func (conf *config) getConfig() *config {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal("yamlFile.Get err:", err)
		return nil
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatal("yamlFile.Get err:", err)
		return nil
	}
	return conf
}

// configInit 获取配置文件
func configInit() {
	var conf config
	conf.getConfig()
	Host = conf.MySQL.Host
	Pass = conf.MySQL.Pass
	Port = conf.MySQL.Port
	User = conf.MySQL.User
	Name = conf.MySQL.Name
	mysql.Init(Host, Port, Name, User, Pass, "utf8mb4", 500, 500)
	mysql.ReleaseMode()
}

// MySQL 连接数据库
func MySQL() *mysql.MysqlDB {
	mysql.UseCache()
	mdb, err := mysql.GetMysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	return mdb
}
