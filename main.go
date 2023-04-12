/**
 * Created by YuYoung on 2023/3/22
 * Description: 入口文件
 */

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	runAddr := viper.GetString("server.host") + ":" + viper.GetString("server.port")
	err := engine.Run(runAddr)
	if err != nil {
		moduleLogger.Error(err)
		panic(err)
	}
}
