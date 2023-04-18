/**
 * Created by YuYoung on 2023/3/22
 * Description: 入口文件
 */

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/database"
	linkHandler "short_link_sys_web_be/handler/link"
	serverHandler "short_link_sys_web_be/handler/server"
	visitHandler "short_link_sys_web_be/handler/visit"
	"short_link_sys_web_be/link_gen"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/middleware"
	"short_link_sys_web_be/router"
	"syscall"
	"time"
)

func init() {
	conf.Init()
	log.Init()
	database.Init()

	link_gen.Init()
	linkHandler.Init()
	visitHandler.Init()
	serverHandler.Init()
	log.GetLogger().Info("all module has init")
}

func main() {
	moduleLogger := log.GetLogger()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(log.MainLogger.Writer()))
	engine.Use(log.Middleware)
	engine.Use(middleware.CrosMiddleware)
	router.LoadAllRouter(engine)
	runAddr := conf.GlobalConfig.GetString("server.host") + ":" + conf.GlobalConfig.GetString("server.port")
	log.GetLogger().Info("server listening on ", runAddr)
	srv := &http.Server{
		Addr:    runAddr,
		Handler: engine,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			moduleLogger.Error(err)
			panic(err)
		}
	}()

	// 阻塞, 等待结束
	sig := <-sigCh
	moduleLogger.Info("receive signal: ", sig, ", start to exit...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		moduleLogger.Error(err)
	}

	terminate()
}

// 资源清理与需要结束的操作
func terminate() {
	link_gen.Terminate()
}
