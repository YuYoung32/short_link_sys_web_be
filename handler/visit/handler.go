/**
 * Created by YuYoung on 2023/3/22
 * Description: 访问相关的 handler
 */

package visit

import (
	_ "embed"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/utils"
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

func testDayIPDataGenerator() []IPSourceResponse {
	provinceList := getRandArr()
	amountList := getRandArr()

	var ipSource []IPSourceResponse
	for i := 0; i < 10; i++ {
		ipSource = append(ipSource, IPSourceResponse{
			Region: ProvinceList[provinceList[i]],
			Amount: amountList[i] + 1,
		})
	}
	return ipSource
}

func testBetweenIPDataGenerator() []IPSourceResponse {
	provinceList := getRandArr()
	amountList := getRandArr()

	var ipSource []IPSourceResponse
	for i := 0; i <= 10; i++ {
		ipSource = append(ipSource, IPSourceResponse{
			Region: ProvinceList[provinceList[i]],
			Amount: amountList[i] + 1,
		})
	}
	return ipSource
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

/*
IPListHandler
IP地址查询方法

	https://ip.taobao.com/instructions
	http://ip.taobao.com/outGetIpInfo?ips=8.8.8.8&accessKey=alibaba-inc
*/
func IPListHandler(ctx *gin.Context) {
	ctx.Set("module", "ip_list_handler")

	begin, end, err := utils.ConvertAndCheckTimeGroup(ctx.Query("begin"), ctx.Query("end"))
	if err != nil {
		common.ErrInvalidArgsResp(ctx)
		return
	}
	if begin == end {
		ctx.JSON(http.StatusOK, testDayIPDataGenerator())
	} else {
		ctx.JSON(http.StatusOK, testBetweenIPDataGenerator())
	}
}
