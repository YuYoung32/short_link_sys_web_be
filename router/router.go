/**
 * Created by YuYoung on 2023/3/22
 * Description: 总路由
 */

package router

import (
	"github.com/gin-gonic/gin"
	"short_link_sys_web_be/handler/auth"
	"short_link_sys_web_be/handler/link"
	"short_link_sys_web_be/handler/server"
	"short_link_sys_web_be/handler/visit"
	"short_link_sys_web_be/middleware"
)

func ServerRouter(engine *gin.Engine) {
	group := engine.Group("/server", middleware.AuthMiddleware)
	group.GET("/info1S", server.RealtimeDataHandler)
	group.GET("/info1Min", server.Info1MinListHandler)
	group.GET("/staticInfo", server.StaticInfoHandler)
}

func VisitRouter(engine *gin.Engine) {
	group := engine.Group("/visit", middleware.AuthMiddleware)
	group.GET("/statics", visit.StaticsListHandler)
	group.GET("/amountTotal", visit.AmountTotalHandler)
	group.GET("/ip", visit.IPListHandler)
	group.POST("/details", visit.DetailsListHandler)
}

func LinkRouter(engine *gin.Engine) {
	group := engine.Group("/link", middleware.AuthMiddleware)
	group.GET("/details", link.DetailsListHandler)
	group.POST("/add", link.AddLinkHandler)
	group.POST("/del", link.DelLinkHandler)
	group.POST("/update", link.UpdateLinkHandler)
	group.GET("/amountTotal", link.AmountTotalHandler)
}

func AuthRouter(engine *gin.Engine) {
	group := engine.Group("/auth")
	group.POST("/login", auth.LoginHandler)
	group.POST("/modify", auth.ChangePasswordHandler)
}

func LoadAllRouter(engine *gin.Engine) {
	ServerRouter(engine)
	VisitRouter(engine)
	LinkRouter(engine)
	AuthRouter(engine)
}
