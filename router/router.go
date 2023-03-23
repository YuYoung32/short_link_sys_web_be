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
	group.GET("/info1S", server.RealtimeDataHandler)
	group.GET("/info1Min", server.Info1MinListHandler)
	group.GET("/staticInfo", server.StaticInfoHandler)
	group.GET("/infoXhr", server.InfoXhrListHandler)
}

func VisitRouter(engine *gin.Engine) {
	group := engine.Group("/visit")
	group.GET("/amountXhr", visit.AmountXHourListHandler)
	group.GET("/amountXhrTotal", visit.AmountXHourTotalHandler)
	group.GET("/ipXhr", visit.IPXHourListHandler)
	group.GET("/details", visit.DetailsListHandler)
}

func LinkRouter(engine *gin.Engine) {
	group := engine.Group("/link")
	group.GET("/details", link.DetailsListHandler)
	group.POST("/add", link.AddLinkHandler)
	group.GET("/del", link.DelLinkHandler)
	group.GET("/amountTotal", link.AmountTotalHandler)
}

func LoadAllRouter(engine *gin.Engine) {
	ServerRouter(engine)
	VisitRouter(engine)
	LinkRouter(engine)
}