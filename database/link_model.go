/**
 * Created by YuYoung on 2023/4/5
 * Description: 链接 ORM Model
 */

package database

import (
	"gorm.io/gorm"
	"short_link_sys_web_be/log"
)

type Link struct {
	ShortLink  string `json:"shortLink" gorm:"primaryKey"`
	LongLink   string `json:"longLink"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime"`
	Comment    string `json:"comment"`
}

func autoMigrateLinkModel(db *gorm.DB) {
	err := db.AutoMigrate(&Link{})
	if err != nil {
		log.MainLogger.WithField("module", "database").Error("auto migrate failed: " + err.Error())
	}
}
