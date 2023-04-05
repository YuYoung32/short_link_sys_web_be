/**
 * Created by YuYoung on 2023/4/5
 * Description: 链接 ORM Model
 */

package database

import (
	"gorm.io/gorm"
	"math/rand"
	"short_link_sys_web_be/log"
	"strconv"
	"time"
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func testLinkDataGenerator() []Link {
	var detailsStore []Link
	for i := 0; i < 100; i++ {
		detailsStore = append(detailsStore, Link{
			ShortLink:  RandomString(5) + strconv.Itoa(i),
			LongLink:   "https://baidu.com",
			CreateTime: time.Now().Unix(),
			Comment:    "无",
		})
	}
	return detailsStore
}
