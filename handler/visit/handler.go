/**
 * Created by YuYoung on 2023/3/22
 * Description: 访问相关的 handler
 */

package visit

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"short_link_sys_web_be/database"
	"short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/utils"
	"time"
)

func isSameDay(timestamp1, timestamp2 int64) bool {
	t1 := time.Unix(timestamp1, 0)
	t2 := time.Unix(timestamp2, 0)
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// sameDayStatics 同一天, 查询当天24小时
func sameDayStatics(whereTemplate string, whereArgs []interface{}, db *gorm.DB) StaticsListResponse {
	logger := log.GetLogger()
	var statics StaticsListResponse
	statics.VisitAmountList = make([]int, 24)
	statics.IPAmountList = make([]int, 24)

	rows, err := db.Model(&database.LinkVisit{}).
		Select(`date_format(from_unixtime(visit_time), '%H') as hour,
					count(*)                     as visitors,
                    count(distinct ip)           as unique_visitors`).
		Where(whereTemplate, whereArgs...).Group("hour").
		Rows()
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(err)
		}
	}(rows)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
	}

	for rows.Next() {
		var hour int
		var visitors int
		var uniqueVisitors int
		if err = rows.Scan(&hour, &visitors, &uniqueVisitors); err != nil {
			logger.Error(err)
			continue
		}
		statics.VisitAmountList[hour] = visitors
		statics.IPAmountList[hour] = uniqueVisitors
	}

	return statics
}

// manyDayStatics 不同天, 查询每天的访问量
func manyDayStatics(whereTemplate string, whereArgs []interface{}, db *gorm.DB) StaticsListResponse {
	logger := log.GetLogger()

	// 生成日期列表, 便于记录当天访问为0的情况
	begin, _ := whereArgs[0].(int64)
	end, _ := whereArgs[1].(int64)
	dateMap := make(map[string]int)
	dateCount := -1
	for ts := begin; ts <= end; ts += 24 * 60 * 60 {
		dateCount++
		date := time.Unix(ts, 0).Format("2006-01-02")
		dateMap[date] = dateCount
	}

	var statics StaticsListResponse
	statics.VisitAmountList = make([]int, dateCount+1)
	statics.IPAmountList = make([]int, dateCount+1)
	rows, err := db.Model(&database.LinkVisit{}).
		Select(`date_format(from_unixtime(visit_time), '%Y-%m-%d') as day,
					count(*)                     as visitors,
					count(distinct ip)           as unique_visitors`).
		Where(whereTemplate, whereArgs...).Group("day").
		Rows()
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(err)
		}
	}(rows)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
	}

	for rows.Next() {
		var day string
		var visitors int
		var uniqueVisitors int
		if err = rows.Scan(&day, &visitors, &uniqueVisitors); err != nil {
			logger.Error(err)
			continue
		}
		statics.VisitAmountList[dateMap[day]] = visitors
		statics.IPAmountList[dateMap[day]] = uniqueVisitors
	}

	return statics
}

// StaticsListHandler 获取指定时间段的访问量, 同天返回24小时, 其余分天
func StaticsListHandler(ctx *gin.Context) {
	begin, end, err := utils.ConvertAndCheckTimeGroup(ctx.Query("begin"), ctx.Query("end"))
	if err != nil {
		common.ErrInvalidArgsResp(ctx)
		return
	}
	shortLink := ctx.Query("shortLink")

	var whereTemplate = "visit_time >= ? and visit_time <= ?"
	var whereArgs = []interface{}{begin, end}
	db := database.GetMysqlInstance()

	if shortLink != "" {
		whereTemplate += " and short_link = ?"
		whereArgs = append(whereArgs, shortLink)
	}
	if isSameDay(begin, end) {
		ctx.JSON(http.StatusOK, sameDayStatics(whereTemplate, whereArgs, db))
	} else {
		ctx.JSON(http.StatusOK, manyDayStatics(whereTemplate, whereArgs, db))
	}
}

// AmountTotalHandler 获取指定时间段的访问量总和
func AmountTotalHandler(ctx *gin.Context) {
	ctx.Set("module", "amount_total_handler")
	begin, end, err := utils.ConvertAndCheckTimeGroup(ctx.Query("begin"), ctx.Query("end"))
	if err != nil {
		common.ErrInvalidArgsResp(ctx)
		return
	}
	ctx.JSON(http.StatusOK, AmountResponse{
		Amount: int(end - begin),
	})
}

func IPListHandler(ctx *gin.Context) {
	logger := log.GetLogger()
	db := database.GetMysqlInstance()
	var ipSourceResponse IPSourceResponse
	begin, end, err := utils.ConvertAndCheckTimeGroup(ctx.Query("begin"), ctx.Query("end"))
	if err != nil {
		common.ErrInvalidArgsResp(ctx)
		return
	}
	shortLink := ctx.Query("shortLink")

	var whereTemplate = "visit_time >= ? and visit_time <= ?"
	var whereArgs = []interface{}{begin, end}
	if shortLink != "" {
		whereTemplate += " and short_link = ?"
		whereArgs = append(whereArgs, shortLink)
	}
	rows, err := db.Model(&database.LinkVisit{}).
		Select("region, count(*) as amount").
		Where(whereTemplate, whereArgs...).Group("region").
		Rows()
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(err)
		}
	}(rows)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
	}

	for rows.Next() {
		var region string
		var amount int
		if err = rows.Scan(&region, &amount); err != nil {
			logger.Error(err)
			continue
		}
		ipSourceResponse.Amount = append(ipSourceResponse.Amount, amount)
		ipSourceResponse.Region = append(ipSourceResponse.Region, region)
	}

	ctx.JSON(http.StatusOK, ipSourceResponse)
}
