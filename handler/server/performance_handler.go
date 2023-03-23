/**
 * Created by YuYoung on 2023/3/22
 * Description: server性能监控handler
 */

package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "short_link_sys_web_be/handler/common"
	"strconv"
)

var testStaticInfoData = StaticInfo{
	MemTotalSize:  10240,
	DiskTotalSize: 13240,
}

func testInfo1MinDataGenerator() Info1Min {
	var info1Min Info1Min
	for i := 0; i < 60; i++ {
		info1Min.CPUUsageRatioSec[i] = i
		info1Min.MemUsageRatioSec[i] = i
		info1Min.DiskUsageRatioSec[i] = i
	}
	return info1Min
}
func testInfoXhrDataGenerator() InfoXhr {
	var infoXhr InfoXhr
	for i := 0; i < 60; i++ {
		infoXhr.TimePoints[i] = i
		infoXhr.CPUUsageRatioMin[i] = i
		infoXhr.MemUsageRatioMin[i] = i
		infoXhr.DiskUsageRatioMin[i] = i
		infoXhr.TTLMin[i] = i
	}
	return infoXhr
}

func StaticInfoHandler(ctx *gin.Context) {
	ctx.Set("module", "static_info_handler")
	ctx.JSON(http.StatusOK, testStaticInfoData)
}

func Info1MinListHandler(ctx *gin.Context) {
	ctx.Set("module", "info_1_min_list_handler")
	ctx.JSON(http.StatusOK, testInfo1MinDataGenerator())
}

func InfoXhrListHandler(ctx *gin.Context) {
	ctx.Set("module", "info_xhr_list_handler")
	x := ctx.Query("x")
	if x == "" {
		ErrMissArgsResp(ctx)
		return
	}
	intX, err := strconv.Atoi(x)
	if err != nil || intX <= 0 || intX > 7*24 {
		ErrInvalidArgsResp(ctx)
		return
	}
	ctx.Set("module", "info_x_hour_handler")
	ctx.JSON(http.StatusOK, testInfoXhrDataGenerator())
}
