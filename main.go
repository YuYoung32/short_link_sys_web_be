/**
 * Created by YuYoung on 2023/3/22
 * Description: 入口文件
 */

package main

import (
	"github.com/gin-gonic/gin"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/database"
	"short_link_sys_web_be/handler/server"
	"short_link_sys_web_be/handler/visit"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/middleware"
	"short_link_sys_web_be/router"
)

func init() {
	conf.Init()
	log.Init()
	database.Init()
	visit.Init()
	server.Init()
	log.GetLogger().Info("all module has init")
}

func main() {
	moduleLogger := log.GetLogger()

	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(log.MainLogger.Writer()))
	engine.Use(log.Middleware)
	engine.Use(middleware.CrosMiddleware)
	router.LoadAllRouter(engine)

	runAddr := conf.GlobalConfig.GetString("server.host") + ":" + conf.GlobalConfig.GetString("server.port")
	log.GetLogger().Info("server listening on ", runAddr)
	err := engine.Run(runAddr)
	if err != nil {
		moduleLogger.Error(err)
		panic(err)
	}
}
