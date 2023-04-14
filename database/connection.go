/**
 * Created by YuYoung on 2023/4/1
 * Description: 数据库连接
 */

package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"short_link_sys_web_be/conf"
	_ "short_link_sys_web_be/conf"
	"short_link_sys_web_be/log"
	"time"
)

var db *gorm.DB

func Init() {
	var err error
	logger := log.GetLogger()

	var dsn = conf.GlobalConfig.GetString("mysql.username") + ":" +
		conf.GlobalConfig.GetString("mysql.password") + "@tcp(" +
		conf.GlobalConfig.GetString("mysql.host") + ":" +
		conf.GlobalConfig.GetString("mysql.port") + ")/" +
		conf.GlobalConfig.GetString("mysql.database")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("failed to connect database: " + err.Error())
		panic("failed to connect database")
	}
	logger.Info("connect database success")
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get sqlDB: " + err.Error())
		panic("failed to get sqlDB")
	}
	sqlDB.SetMaxIdleConns(conf.GlobalConfig.GetInt("mysql.maxIdleConns"))
	sqlDB.SetMaxOpenConns(conf.GlobalConfig.GetInt("mysql.maxOpenConns"))
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(conf.GlobalConfig.GetInt("mysql.connMaxLifetime")))

	autoMigrate()
}

// GetDBInstance 获取数据库实例, 其他包使用
func GetDBInstance() *gorm.DB {
	if db == nil {
		log.GetLogger().Error("db is nil")
		panic("db is nil")
	}
	return db
}

func autoMigrate() {
	db := GetDBInstance()
	autoMigrateLinkModel(db)
}
