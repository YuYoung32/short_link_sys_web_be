/**
 * Created by YuYoung on 2023/3/22
 * Description: 总路由
 */

package router

import (
	"github.com/gin-gonic/gin"
	"short_link_sys_web_be/handler/link"
	"short_link_sys_web_be/handler/server"
	"short_link_sys_web_be/handler/visit"
)

func ServerRouter(engine *gin.Engine) {
	group := engine.Group("/server")
	group.GET("/info_1s", server.RealtimeDataHandler)
	group.GET("/static_info", server.StaticInfoHandler)
	group.GET("/cpu_xhr", server.CPUUsageRatiosHandler)
	group.GET("/memory_xhr", server.MemoryUsageRatiosHandler)
	group.GET("/disk_xhr", server.DiskUsageRatiosHandler)
	group.GET("/ttl_xhr", server.TTLHandler)
}

func VisitRouter(engine *gin.Engine) {
	group := engine.Group("/visit")
	group.GET("/amount_xhr", visit.AmountXHourHandler)
	group.GET("/amount", visit.AmountHandler)
	group.GET("/ip_xhr", visit.IPXHourHandler)
	group.GET("/details", visit.DetailsHandler)
}

func LinkRouter(engine *gin.Engine) {
	group := engine.Group("/link")
	group.GET("/details", link.GetLinkDetailsListHandler)
	group.GET("/add", link.AddLinkHandler)
	group.GET("/del", link.DelLinkHandler)
	group.GET("/total_amount", link.GetTotalAmountHandler)
}

func LoadAllRouter(engine *gin.Engine) {
	ServerRouter(engine)
	VisitRouter(engine)
	LinkRouter(engine)
}
