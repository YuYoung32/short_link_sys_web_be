/**
 * Created by YuYoung on 2023/4/1
 * Description: 访问 ORM Model
 */

package database

import (
	"gorm.io/gorm"
	"short_link_sys_web_be/log"
)

type Visit struct {
	ShortLink string `json:"shortLink" gorm:"primaryKey"`
	LongLink  string `json:"longLink"`
	Comment   string `json:"comment"`
	IP        string `json:"ip"`
	Region    string `json:"region"`
	VisitTime int64  `json:"visitTime" gorm:"autoCreateTime"`
}

func autoMigrateVisitModel(db *gorm.DB) {
	err := db.AutoMigrate(&Visit{})
	if err != nil {
		log.MainLogger.WithField("module", "database").Error("auto migrate failed: " + err.Error())
	}
}
