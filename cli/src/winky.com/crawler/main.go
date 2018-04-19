package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"os"
)

var (
	conf   config
	pool   *redis.Pool
	engine *xorm.Engine
	imgUrl = [4]string{"http://img6.winky.cn", "http://img7.winky.cn", "http://img8.winky.cn", "http://img9.winky.cn"}
)

const (
	MaxOpenDbConnect int = 50
	MaxIdleConn      int = 5
	MaxOpenGoRoutine     = 20

	BusinessArticle = 1
)

func main() {

	if _, err := toml.DecodeFile("./conf/config.toml", &conf); err != nil {
		log.Fatal("toml decode conf err : ", err.Error())
		return
	}

	file, err := os.OpenFile(conf.Log.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("open log failed: ", err)
	}
	defer file.Close()
	log.SetOutput(file)

	redisAddr := fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port)
	pool = newRedisPool(redisAddr, conf.Redis.Pswd, conf.Redis.Database, conf.Redis.Timeout)
	defer pool.Close()

	mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.Mysql.User,
		conf.Mysql.Pswd,
		conf.Mysql.Host,
		conf.Mysql.Port,
		conf.Mysql.Database,
		conf.Mysql.Charset,
	)
	engine, err = xorm.NewEngine("mysql", mysqlDsn)
	if err != nil {
		log.Panic("open mysql failed : ", err)
	}
	defer engine.Close()

	engine.SetMaxIdleConns(MaxIdleConn)      // 连接池空闲数
	engine.SetMaxOpenConns(MaxOpenDbConnect) // 最大连接数

	engine.ShowSQL(conf.Log.SqlLog)
	engine.SetLogger(xorm.NewSimpleLogger(file))

	err = engine.Ping()
	if err != nil {
		log.Panic("mysql connect failed : ", err)
	}
	log.Println("Init Finished.")

	// 审核文章
	go procBiz()
	select {}

	// 采集
	//procCrawler()
	//article(13020)
}
