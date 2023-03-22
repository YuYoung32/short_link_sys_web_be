/**
 * Created by YuYoung on 2023/3/22
 * Description: 入口文件
 */

package main

import (
	"github.com/gin-gonic/gin"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/middleware"
)

func main() {
	engine := gin.New()
	engine.Use(log.Middleware)
	engine.Use(middleware.CrosMiddleware)
}
