/**
 * Created by YuYoung on 2023/4/7
 * Description: 服务器本身信息
 */

package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var testStaticInfoData = StaticInfo{
	CPUStaticInfo: CPUStaticInfo{
		Name:      "Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
		CoreNum:   4,
		ThreadNum: 8,
		CacheSize: 6,
	},
	MemStaticInfo: MemStaticInfo{
		PhysicalTotalSize: 16 * 1024 * 1024 * 1024,
		SwapTotalSize:     0,
	},
	DiskStaticInfo: DiskStaticInfo{
		DiskTotalSize: 40 * 1024 * 1024 * 1024,
	},
	NetStaticInfo: NetStaticInfo{
		IPv4: "10.1.23.1",
		MAC:  "1F:2A:3b:4C:5D:6E",
	},
}

func StaticInfoHandler(ctx *gin.Context) {
	ctx.Set("module", "static_info_handler")
	ctx.JSON(http.StatusOK, testStaticInfoData)
}
