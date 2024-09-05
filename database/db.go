package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"normal_web/util"
	"normal_web/util/logger"
	"sync"
)

var (
	mysql_db    *gorm.DB
	mysqlDBOnce sync.Once

	redis_client    *redis.Client
	redisClientOnce sync.Once
)

func connectMysql(dbname, host, username, password string, port int) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.NewGormLogger(), PrepareStmt: true})
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
		if mysql_db == nil {
			mysqlConf := util.CreateConfig("mysql")
			dbname := mysqlConf.GetString("database")
			host := mysqlConf.GetString("host")
			username := mysqlConf.GetString("username")
			password := mysqlConf.GetString("password")
			port := mysqlConf.GetInt("port")
			mysql_db = connectMysql(dbname, host, username, password, port)
		}
	})
	return mysql_db
}

func createRedisClient(addr, password string, db int) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	cxt := context.Background()
	if err := cli.Ping(cxt).Err(); err != nil {
		logger.Error(fmt.Sprintf("connect redis %s failed: %s", addr, err))
	} else {
		logger.Info(fmt.Sprintf("connect redis %s successful!", addr))
	}
	return cli
}

func GetRedisClient() *redis.Client {
	redisClientOnce.Do(func() {
		if redis_client != nil {
			redis_client.Close()
		}
		redisConf := util.CreateConfig("redis")
		redis_client = createRedisClient(redisConf.GetString("addr"), redisConf.GetString("password"), redisConf.GetInt("db"))

	})
	return redis_client
}
