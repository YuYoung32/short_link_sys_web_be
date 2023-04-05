/**
 * Created by YuYoung on 2023/4/1
 * Description: 数据库连接
 */

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"short_link_sys_web_be/log"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	logger := log.MainLogger.WithField("module", "database_init")
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		logger.Error("failed to connect database: " + err.Error())
		panic("failed to connect database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get sqlDB: " + err.Error())
		panic("failed to get sqlDB")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	autoMigrate()
	generateTestData()
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
	autoMigrateVisitModel(db)
	autoMigrateLinkModel(db)
}

// TODO 生成测试数据 删除
func generateTestData() {
	db := GetDBInstance()
	visitList := testVisitDataGenerator()
	db.Create(&visitList)
	linkList := testLinkDataGenerator()
	db.Create(&linkList)
}
