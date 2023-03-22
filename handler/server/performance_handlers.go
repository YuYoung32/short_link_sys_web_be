/**
 * Created by YuYoung on 2023/3/22
 * Description: server性能监控handler
 */

package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var testStaticInfo = StaticInfo{
	MemTotalSize:  10240,
	DiskTotalSize: 13240,
}

func StaticInfoHandler(ctx *gin.Context) {
	//ModuleLogger := log.MainLogger.WithField("module", "static_info_handler")
	ctx.JSON(http.StatusOK, testStaticInfo)
}

func CPUUsageRatiosHandler(ctx *gin.Context) {
	//ModuleLogger := log.MainLogger.WithField("module", "cpu_usage_ratios_handler")
}

func MemoryUsageRatiosHandler(ctx *gin.Context) {
	//ModuleLogger := log.MainLogger.WithField("module", "memory_usage_ratios_handler")
}

func DiskUsageRatiosHandler(ctx *gin.Context) {
	//ModuleLogger := log.MainLogger.WithField("module", "disk_usage_ratios_handler")
}

func TTLHandler(ctx *gin.Context) {
	//ModuleLogger := log.MainLogger.WithField("module", "ttl_handler")
}
