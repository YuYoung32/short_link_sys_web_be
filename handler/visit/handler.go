/**
 * Created by YuYoung on 2023/3/22
 * Description:
 */

package visit

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "short_link_sys_web_be/handler/common"
	"strconv"
)

func testAmountListDataGenerator(x int) AmountTime {
	var amountTime AmountTime
	for i := 0; i < x; i++ {
		amountTime.TimePoints = append(amountTime.TimePoints, i)
		amountTime.Amount = append(amountTime.Amount, Amount(i))
	}
	return amountTime
}

func testIPSourceListDataGenerator() []IPSource {
	var ipSourceList []IPSource
	for i := 0; i < 10; i++ {
		ipSourceList = append(ipSourceList, IPSource{
			Region: "浙江",
			Amount: i,
		})
	}
	return ipSourceList
}

func testDetailsListDataGenerator() []Details {
	var detailsList []Details
	for i := 0; i < 10; i++ {
		detailsList = append(detailsList, Details{
			LongUrl:   "https://www.baidu.com/",
			ShortUrl:  "sdvser",
			IP:        "",
			Region:    "浙江",
			OS:        "Windows",
			Timestamp: "1679555628",
		})
	}
	return detailsList
}

func AmountXHourTotalHandler(ctx *gin.Context) {
	ctx.Set("module", "amount_x_hour_total_handler")
	x := ctx.Query("x")
	if x == "" {
		ErrMissArgsResp(ctx)
		return
	}
	intX, err := strconv.Atoi(x)
	if err != nil || intX <= 0 || intX > MaxSearchHours {
		ErrInvalidArgsResp(ctx)
		return
	}
	ctx.JSON(http.StatusOK, AmountTotal{
		Amount: Amount(intX),
	})
}

// AmountXHourListHandler 获取最近x小时的访问量, 分时成list
func AmountXHourListHandler(ctx *gin.Context) {
	ctx.Set("module", "amount_x_hour_list_handler")
	x := ctx.Query("x")
	if x == "" {
		ErrMissArgsResp(ctx)
		return
	}
	intX, err := strconv.Atoi(x)
	if err != nil || intX <= 0 || intX > MaxSearchHours {
		ErrInvalidArgsResp(ctx)
		return
	}
	ctx.JSON(http.StatusOK, testAmountListDataGenerator(intX))
}

/*
IP地址查询方法
	https://ip.taobao.com/instructions
	http://ip.taobao.com/outGetIpInfo?ips=8.8.8.8&accessKey=alibaba-inc
*/

func IPXHourListHandler(ctx *gin.Context) {
	ctx.Set("module", "ip_x_hour_list_handler")
	x := ctx.Query("x")
	if x == "" {
		ErrMissArgsResp(ctx)
		return
	}
	intX, err := strconv.Atoi(x)
	if err != nil || intX <= 0 || intX > MaxSearchHours {
		ErrInvalidArgsResp(ctx)
		return
	}
	ctx.JSON(http.StatusOK, testIPSourceListDataGenerator())
}

func DetailsListHandler(ctx *gin.Context) {
	ctx.Set("module", "details_list_handler")
	amount := ctx.Query("amount")
	if amount == "" {
		ErrMissArgsResp(ctx)
		return
	}
	intAmount, err := strconv.Atoi(amount)
	if err != nil || intAmount <= 0 || intAmount > MaxSearchAmount {
		ErrInvalidArgsResp(ctx)
		return
	}
	ctx.JSON(http.StatusOK, testDetailsListDataGenerator())
}
