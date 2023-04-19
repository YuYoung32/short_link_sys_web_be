/**
 * Created by YuYoung on 2023/4/4
 * Description: 详情列表, 处理函数过程故单独开一个文件
 */

package visit

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"short_link_sys_web_be/database"
	. "short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/utils"
	"strconv"
	"strings"
)

// 检查传入的IP地址是否合法, 并将其中的通配符替换为%, 用于数据库查询
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
		var begin, end int64
		begin, end, err := utils.ConvertAndCheckTimeGroup(queryDetailsBind.RangeTime[0], queryDetailsBind.RangeTime[1])
		if err != nil {
			ErrInvalidArgsResp(ctx)
			return
		}
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
