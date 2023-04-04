/**
 * Created by YuYoung on 2023/3/22
 * Description: 访问相关的 handler
 */

package visit

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"short_link_sys_web_be/database"
	. "short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/log"
	"strconv"
	"strings"
	"time"
)

const OneDaySub1Sec = 24*time.Hour - time.Second

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

func testDayVADataGenerator(day time.Time) AmountTimeResponse {
	var amountTime AmountTimeResponse
	for i := 0; i < 24; i++ {
		amountTime.Amount = append(amountTime.Amount, day.Day()+i)
	}
	return amountTime
}

func testBetweenVADataGenerator(begin time.Time, end time.Time) AmountTimeResponse {
	var amountTime AmountTimeResponse
	for i := 0; i <= int(end.Sub(begin).Hours()/24); i++ {
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

func testDayIPDataGenerator(day time.Time) []IPSourceResponse {
	provinceList := getRandArr()
	amountList := getRandArr()

	var ipSource []IPSourceResponse
	for i := 0; i < day.Day(); i++ {
		ipSource = append(ipSource, IPSourceResponse{
			Region: ProvinceList[provinceList[i]],
			Amount: amountList[i] + 1,
		})
	}
	return ipSource
}

func testBetweenIPDataGenerator(begin time.Time, end time.Time) []IPSourceResponse {
	provinceList := getRandArr()
	amountList := getRandArr()

	var ipSource []IPSourceResponse
	for i := begin.Day(); i <= end.Day(); i++ {
		ipSource = append(ipSource, IPSourceResponse{
			Region: ProvinceList[provinceList[i]],
			Amount: amountList[i] + 1,
		})
	}
	return ipSource
}

// parseTimeOrBadResponse 解析时间参数begin和end 如果参数错误(任意一个为空或begin>end或end>=今天)则返回错误和响应HTTP 400
func parseTimeOrBadResponse(ctx *gin.Context) (error, time.Time, time.Time) {
	begin := ctx.Query("begin")
	end := ctx.Query("end")
	if begin == "" || end == "" {
		ErrMissArgsResp(ctx)
		return errors.New(""), time.Time{}, time.Time{}
	}
	beginTime, beginErr := time.Parse("20060102", begin)
	endTime, endErr := time.Parse("20060102", end)
	if beginErr != nil || endErr != nil || beginTime.After(endTime) || endTime.After(time.Now()) {
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

// AmountTotalHandler 获取指定时间段的访问量总和
func AmountTotalHandler(ctx *gin.Context) {
	ctx.Set("module", "amount_total_handler")
	begin := ctx.Query("begin")
	end := ctx.Query("end")
	var beginTimestamp, endTimestamp int64
	yesterday := time.Now().AddDate(0, 0, -1)
	endTimestamp = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, time.Local).Unix()
	if begin != "" {
		if beginTime, beginErr := time.Parse("20060102", begin); beginErr == nil {
			beginTimestamp = beginTime.Unix()
		} else {
			ErrInvalidArgsResp(ctx)
			return
		}
	}
	if end != "" {
		if endTime, endErr := time.Parse("20060102", end); endErr == nil {
			tmpEndTimestamp := time.Unix(endTime.Unix(), 0).Add(OneDaySub1Sec).Unix()
			if tmpEndTimestamp > endTimestamp {
				ErrInvalidArgsResp(ctx)
				return
			}
			endTimestamp = tmpEndTimestamp
		} else {
			ErrInvalidArgsResp(ctx)
			return
		}
	}
	if beginTimestamp >= endTimestamp {
		ErrInvalidArgsResp(ctx)
		return
	}
	ctx.JSON(http.StatusOK, AmountResponse{
		Amount: int(endTimestamp - beginTimestamp),
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

func checkAndReplaceIPStr(ip []string) ([]string, error) {
	for i, item := range ip {
		ipNums := strings.Split(item, ".")
		if len(ipNums) != 4 {
			return []string{}, errors.New("invalid ip address")
		}
		for j, ipNum := range ipNums {
			if num, err := strconv.Atoi(ipNum); err != nil {
				ipNums[j] = "%"
			} else {
				if num < 0 || num > 255 {
					return []string{}, errors.New("invalid ip address")
				}
			}
		}
		ip[i] = strings.Join(ipNums, ".")
	}
	return ip, nil
}

func DetailsListHandler(ctx *gin.Context) {
	ctx.Set("module", "details_list_handler")
	db := database.GetDBInstance()

	//region 获取POST body参数
	type QueryDetailsBind struct {
		ShortLink []string `json:"shortLink"`
		LongLink  []string `json:"longLink"`
		Comment   []string `json:"comment"`
		Region    []string `json:"region"`
		IP        []string `json:"ip"`
		RangeTime []string `json:"rangeTime"`
	}
	queryDetailsBind := QueryDetailsBind{}
	if err := ctx.ShouldBind(&queryDetailsBind); err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}
	//endregion

	//region 构造查询条件
	// 内部为同一字段的OR查询, 外部为不同字段的AND查询
	var queryTemplateList []string
	var queryArgsList []interface{}

	// 一个简单的封装, 构造like查询, 多个like之间使用OR连接, 外部加上()
	buildLikeQuery := func(args []string, str string) {
		if len(args) < 1 {
			return
		}
		var queryTemplateOr []string
		for _, arg := range args {
			queryArgsList = append(queryArgsList, "%"+arg+"%")
			queryTemplateOr = append(queryTemplateOr, str)
		}

		queryTemplateList = append(queryTemplateList, "("+strings.Join(queryTemplateOr, " or ")+")")
	}
	buildLikeQuery(queryDetailsBind.ShortLink, "short_link like (?)")
	buildLikeQuery(queryDetailsBind.LongLink, "long_link like (?)")
	buildLikeQuery(queryDetailsBind.Comment, "comment like (?)")

	if len(queryDetailsBind.Region) > 0 {
		queryTemplateList = append(queryTemplateList, "region in (?)")
		queryArgsList = append(queryArgsList, queryDetailsBind.Region)
	}

	if newIPs, err := checkAndReplaceIPStr(queryDetailsBind.IP); err == nil && len(newIPs) > 0 {
		var queryTemplateOr []string
		for _, ip := range newIPs {
			queryArgsList = append(queryArgsList, ip)
			queryTemplateOr = append(queryTemplateOr, "ip like (?)")
		}
		queryTemplateList = append(queryTemplateList, "("+strings.Join(queryTemplateOr, " or ")+")")
	} else if err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}

	if len(queryDetailsBind.RangeTime) == 2 {
		beginTime, err := time.Parse("20060102", queryDetailsBind.RangeTime[0])
		endTime, err := time.Parse("20060102", queryDetailsBind.RangeTime[1])
		if err != nil {
			ErrInvalidArgsResp(ctx)
			return
		}
		begin := beginTime.Unix()
		end := time.Unix(endTime.Unix(), 0).Add(OneDaySub1Sec).Unix()
		queryArgsList = append(queryArgsList, begin, end)
		queryTemplateList = append(queryTemplateList, "visit_time between (?) and (?)")
	} else if len(queryDetailsBind.RangeTime) != 0 {
		ErrInvalidArgsResp(ctx)
		return
	}
	//endregion

	var resp DetailsListResponse
	db.Where(strings.Join(queryTemplateList, " and "), queryArgsList...).Find(&resp.VisitDetails).Count(&resp.VisitDetailsAmount)
	ctx.JSON(http.StatusOK, resp)
}
