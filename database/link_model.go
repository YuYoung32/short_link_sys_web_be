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
	ShortLink  string `json:"shortLink" gorm:"type:varchar(255) COLLATE utf8_bin;primaryKey"`
	LongLink   string `json:"longLink"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime int64  `json:"updateTime" gorm:"autoUpdateTime"`
	Comment    string `json:"comment"`
}

func autoMigrateLinkModel(db *gorm.DB) {
	err := db.AutoMigrate(&Link{})
	if err != nil {
		log.GetLogger().Error("auto migrate failed: " + err.Error())
	}
}
