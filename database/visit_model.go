/**
 * Created by YuYoung on 2023/4/1
 * Description: 访问 ORM Model
 */

package database

import (
	"gorm.io/gorm"
	"short_link_sys_web_be/log"
	"strconv"
	"time"
)

type Details struct {
	ShortLink string `json:"shortLink" gorm:"primaryKey"`
	LongLink  string `json:"longLink"`
	Comment   string `json:"comment"`
	IP        string `json:"ip"`
	Region    string `json:"region"`
	VisitTime int64  `json:"visitTime" gorm:"autoCreateTime"`
}

func autoMigrateVisitModel(db *gorm.DB) {
	err := db.AutoMigrate(&Details{})
	if err != nil {
		log.MainLogger.WithField("module", "database").Error("auto migrate failed: " + err.Error())
	}
}

func testDetailsDataGenerator() []Details {
	var detailsList []Details
	for i := 0; i < 10; i++ {
		detailsList = append(detailsList, Details{
			ShortLink: "test" + strconv.Itoa(i),
			LongLink:  "test" + strconv.Itoa(i),
			Comment:   "test" + strconv.Itoa(i),
			IP:        "192.168.1.12",
			Region:    "浙江",
			VisitTime: time.Now().Unix(),
		})
	}
	return detailsList
}
