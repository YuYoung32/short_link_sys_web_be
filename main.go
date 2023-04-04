/**
 * Created by YuYoung on 2023/3/22
 * Description: 入口文件
 */

package main

import (
	"github.com/gin-gonic/gin"
	_ "short_link_sys_web_be/database"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/middleware"
	"short_link_sys_web_be/router"
)

func main() {
	moduleLogger := log.MainLogger.WithField("module", "main")

	engine := gin.New()
	engine.Use(log.Middleware)
	engine.Use(middleware.CrosMiddleware)
	router.LoadAllRouter(engine)

	err := engine.Run(":8081")
	if err != nil {
		moduleLogger.Error(err)
		panic(err)
	}
}
