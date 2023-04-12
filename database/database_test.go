/**
 * Created by YuYoung on 2023/4/12
 * Description:
 */

package database

import (
	"strconv"
	"testing"
	"time"
)

func linkDataGenerator() []Link {
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

func visitDataGenerator() []Visit {
	var detailsList []Visit
	for i := 0; i < 10; i++ {
		detailsList = append(detailsList, Visit{
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

func TestGenerateData(t *testing.T) {
	db := GetDBInstance()
	linkData := linkDataGenerator()
	db.Create(&linkData)
	visitData := visitDataGenerator()
	db.Create(&visitData)
}
