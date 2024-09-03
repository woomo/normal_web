package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"normal_web/util"
	"normal_web/util/logger"
	"sync"
)

var (
	db          *gorm.DB
	mysqlDBOnce sync.Once
)

func connectMysql(dbname, host, username, password string, port int) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.NewGormLogger(), PrepareStmt: true})
	if err != nil {
		msg := fmt.Sprintf("connect to mysql use dsn %s failed: %s", dsn, err)
		logger.Error(msg)
		panic(errors.New(msg))
	}

	//设置数据库连接池参数，提高并发性能
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
	return db
}

func GetMysqlDB() *gorm.DB {
	mysqlDBOnce.Do(func() {
		if db == nil {
			mysqlConf := util.CreateConfig("mysql")
			dbname := mysqlConf.GetString("database")
			host := mysqlConf.GetString("host")
			username := mysqlConf.GetString("username")
			password := mysqlConf.GetString("password")
			port := mysqlConf.GetInt("port")
			connectMysql(dbname, host, username, password, port)
		}
	})
	return db
}
