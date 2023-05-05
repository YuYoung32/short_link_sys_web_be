/**
 * Created by YuYoung on 2023/4/1
 * Description: 数据库连接
 */

package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/log"
	"time"
)

var mysqlDB *gorm.DB
var redisDB *redis.Client

func init() {
	var err error
	moduleLogger := log.GetLogger()

	var dsn = conf.GlobalConfig.GetString("mysql.username") + ":" +
		conf.GlobalConfig.GetString("mysql.password") + "@tcp(" +
		conf.GlobalConfig.GetString("mysql.host") + ":" +
		conf.GlobalConfig.GetString("mysql.port") + ")/" +
		conf.GlobalConfig.GetString("mysql.database")
	mysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(log.MainLogger, logger.Config{
			SlowThreshold:             0,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Warn,
		}),
	})
	if err != nil {
		moduleLogger.Error("failed to connect database: " + err.Error())
		panic("failed to connect database")
	}
	moduleLogger.Info("connect database success")
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		moduleLogger.Error("failed to get sqlDB: " + err.Error())
		panic("failed to get sqlDB")
	}
	sqlDB.SetMaxIdleConns(conf.GlobalConfig.GetInt("mysql.maxIdleConns"))
	sqlDB.SetMaxOpenConns(conf.GlobalConfig.GetInt("mysql.maxOpenConns"))
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(conf.GlobalConfig.GetInt("mysql.connMaxLifetime")))

	autoMigrate()

	// redis初始化
	redisDB = redis.NewClient(&redis.Options{
		Addr:     conf.GlobalConfig.GetString("redis.host") + ":" + conf.GlobalConfig.GetString("redis.port"),
		Password: conf.GlobalConfig.GetString("redis.password"),
		DB:       conf.GlobalConfig.GetInt("redis.db"),
	})
	_, err = redisDB.Ping(context.Background()).Result()
	if err != nil {
		moduleLogger.Error("failed to connect redis: " + err.Error())
		panic("failed to connect redis")
	}
	moduleLogger.Info("connect redis success")
}

// GetMysqlInstance 获取数据库实例, 其他包使用
func GetMysqlInstance() *gorm.DB {
	if mysqlDB == nil {
		log.GetLogger().Error("mysqlDB is nil")
		panic("mysqlDB is nil")
	}
	return mysqlDB
}

// GetRedisInstance 获取redis实例, 其他包使用
func GetRedisInstance() *redis.Client {
	if redisDB == nil {
		log.GetLogger().Error("rdb is nil")
	}
	return redisDB
}

func autoMigrate() {
	db := GetMysqlInstance()
	autoMigrateLinkModel(db)
}
