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
	"sync"
)

type info1MinMutex struct {
	mutex    sync.Mutex
	info1Min Info1Min
}

var (
	testStaticInfoData = StaticInfo{
		MemTotalSize:  10240,
		DiskTotalSize: 13240,
	}
	info info1MinMutex
)

func init() {
	go func() {
		for data := range realTimeDataTransfer {
			info.mutex.Lock()
			pushAndPopArr(&info.info1Min, data)
			info.mutex.Unlock()
		}
	}()
}

func pushAndPopArr(info1Min *Info1Min, data Info1s) {
	for i := 0; i < 59; i++ {
		info1Min.CPUUsageRatioSec[i] = info1Min.CPUUsageRatioSec[i+1]
		info1Min.MemUsageRatioSec[i] = info1Min.MemUsageRatioSec[i+1]
		info1Min.DiskUsageRatioSec[i] = info1Min.DiskUsageRatioSec[i+1]
	}
	info1Min.CPUUsageRatioSec[59] = data.CPUUsageRatioSec
	info1Min.MemUsageRatioSec[59] = data.MemUsageRatioSec
	info1Min.DiskUsageRatioSec[59] = data.DiskUsageRatioSec
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
	ctx.JSON(http.StatusOK, info.info1Min)
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
