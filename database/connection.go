/**
 * Created by YuYoung on 2023/4/1
 * Description: 数据库连接
 */

package database

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "short_link_sys_web_be/conf"
	"short_link_sys_web_be/log"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	logger := log.MainLogger.WithField("module", "database_init")

	var dsn = viper.GetString("mysql.username") + ":" +
		viper.GetString("mysql.password") + "@tcp(" +
		viper.GetString("mysql.host") + ":" +
		viper.GetString("mysql.port") + ")/" +
		viper.GetString("mysql.database")
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
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	autoMigrate()
}

// GetDBInstance 获取数据库实例, 其他包使用
func GetDBInstance() *gorm.DB {
	if db == nil {
		log.MainLogger.WithField("module", "database").Error("db is nil")
		panic("db is nil")
	}
	return db
}

func autoMigrate() {
	db := GetDBInstance()
	autoMigrateLinkModel(db)
}
