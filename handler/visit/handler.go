/**
 * Created by YuYoung on 2023/3/22
 * Description:
 */

package visit

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	. "short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/log"
	"strconv"
	"time"
)

var (
	//go:embed province_code.json
	provinceCodeFile []byte
	ProvinceList     []string
	CodeList         []string
	ProvinceToCode   = make(map[string]string)
	CodeToProvince   = make(map[string]string)
)

func init() {
	logger := log.MainLogger.WithField("module", "visit_handler_init")

	var data []map[string]string
	err := json.Unmarshal(provinceCodeFile, &data)
	if err != nil {
		logger.Error("Unmarshal province_code.json failed: " + err.Error())
		return
	}
	for _, item := range data {
		ProvinceList = append(ProvinceList, item["name"])
		CodeList = append(CodeList, item["code"])
		ProvinceToCode[item["name"]] = item["code"]
		CodeToProvince[item["code"]] = item["name"]
	}
}

func testDayVADataGenerator(day time.Time) AmountTime {
	var amountTime AmountTime
	for i := 0; i < 24; i++ {
		amountTime.Amount = append(amountTime.Amount, day.Day()+i)
	}
	return amountTime
}

func testBetweenVADataGenerator(begin time.Time, end time.Time) AmountTime {
	var amountTime AmountTime
	for i := begin.Day(); i <= end.Day(); i++ {
		amountTime.Amount = append(amountTime.Amount, i)
	}
	return amountTime
}

// getRandArr 生成随机数组 数组内容0-32
func getRandArr() []int {
	arr := make([]int, 33)
	for i := 0; i < len(arr); i++ {
		arr[i] = i
	}
	rand.Seed(time.Now().UnixNano())

	// 使用 Fisher-Yates 洗牌算法打乱数组
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func testDayIPDataGenerator(day time.Time) []IPSource {
	provinceList := getRandArr()
	amountList := getRandArr()

	var ipSource []IPSource
	for i := 0; i < day.Day(); i++ {
		ipSource = append(ipSource, IPSource{
			Region: ProvinceList[provinceList[i]],
			Amount: amountList[i] + 1,
		})
	}
	return ipSource
}

func testBetweenIPDataGenerator(begin time.Time, end time.Time) []IPSource {
	provinceList := getRandArr()
	amountList := getRandArr()

	var ipSource []IPSource
	for i := begin.Day(); i <= end.Day(); i++ {
		ipSource = append(ipSource, IPSource{
			Region: ProvinceList[provinceList[i]],
			Amount: amountList[i] + 1,
		})
	}
	return ipSource
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

func parseTimeOrBadResponse(ctx *gin.Context) (error, time.Time, time.Time) {
	begin := ctx.Query("begin")
	end := ctx.Query("end")
	beginTime, err := time.Parse("20060102", begin)
	endTime, err := time.Parse("20060102", end)
	if begin == "" || end == "" {
		ErrMissArgsResp(ctx)
		return errors.New(""), beginTime, endTime
	}
	if err != nil || beginTime.After(endTime) || endTime.Day() >= time.Now().Day() {
		ErrInvalidArgsResp(ctx)
		return errors.New(""), beginTime, endTime
	}
	return nil, beginTime, endTime
}

// AmountListHandler 获取指定时间段的访问量, 同天返回24小时, 其余分天
func AmountListHandler(ctx *gin.Context) {
	ctx.Set("module", "amount_list_handler")

	err, beginTime, endTime := parseTimeOrBadResponse(ctx)
	if err != nil {
		return
	}
	if beginTime.Equal(endTime) {
		ctx.JSON(http.StatusOK, testDayVADataGenerator(beginTime))
	} else {
		ctx.JSON(http.StatusOK, testBetweenVADataGenerator(beginTime, endTime))
	}
}

/*
IP地址查询方法
	https://ip.taobao.com/instructions
	http://ip.taobao.com/outGetIpInfo?ips=8.8.8.8&accessKey=alibaba-inc
*/

func IPListHandler(ctx *gin.Context) {
	ctx.Set("module", "ip_list_handler")

	err, beginTime, endTime := parseTimeOrBadResponse(ctx)
	if err != nil {
		return
	}
	if beginTime.Equal(endTime) {
		ctx.JSON(http.StatusOK, testDayIPDataGenerator(beginTime))
	} else {
		ctx.JSON(http.StatusOK, testBetweenIPDataGenerator(beginTime, endTime))
	}
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
