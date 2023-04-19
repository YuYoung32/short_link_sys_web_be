/**
 * Created by YuYoung on 2023/3/22
 * Description: 访问相关的 handler
 */

package visit

import (
	"database/sql"
	_ "embed"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"short_link_sys_web_be/database"
	"short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/utils"
	"time"
)

func Init() {
}

// 生成24小时的
func testDayVADataGenerator(begin int64) StaticsListResponse {
	var statics StaticsListResponse
	for i := 0; i < 24; i++ {
		statics.VisitAmountList = append(statics.VisitAmountList, begin/1000+int64(i))
		statics.IPAmountList = append(statics.IPAmountList, begin/1000+10*int64(i))
	}
	return statics
}

// 按天生成
func testBetweenVADataGenerator(begin int64, end int64) StaticsListResponse {
	var statics StaticsListResponse
	for i := 0; i <= int(time.Unix(end, 0).Sub(time.Unix(begin, 0)).Hours()/24); i++ {
		statics.VisitAmountList = append(statics.VisitAmountList, begin/1000+int64(i))
		statics.IPAmountList = append(statics.IPAmountList, begin/1000+10*int64(i))
	}
	return statics
}

func isSameDay(timestamp1, timestamp2 int64) bool {
	t1 := time.Unix(timestamp1, 0)
	t2 := time.Unix(timestamp2, 0)
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// StaticsListHandler 获取指定时间段的访问量, 同天返回24小时, 其余分天
func StaticsListHandler(ctx *gin.Context) {
	ctx.Set("module", "amount_list_handler")
	begin, end, err := utils.ConvertAndCheckTimeGroup(ctx.Query("begin"), ctx.Query("end"))
	if err != nil {
		common.ErrInvalidArgsResp(ctx)
		return
	}
	if isSameDay(begin, end) {
		ctx.JSON(http.StatusOK, testDayVADataGenerator(begin))
	} else {
		ctx.JSON(http.StatusOK, testBetweenVADataGenerator(begin, end))
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
	db := database.GetDBInstance()
	var ipSourceResponse IPSourceResponse
	begin, end, err := utils.ConvertAndCheckTimeGroup(ctx.Query("begin"), ctx.Query("end"))

	rows, err := db.Model(&database.LinkVisit{}).
		Select("ip, count(*) as amount").
		Where("visit_time >= ? and visit_time <= ?", begin, end).Group("ip").
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
		var ip string
		var amount int
		if err = rows.Scan(&ip, &amount); err != nil {
			logger.Error(err)
			continue
		}
		ipSourceResponse.Amount = append(ipSourceResponse.Amount, amount)
		ipSourceResponse.Region = append(ipSourceResponse.Region, ip)
	}

	if err != nil {
		common.ErrInvalidArgsResp(ctx)
		return
	}
	ctx.JSON(http.StatusOK, ipSourceResponse)
}
