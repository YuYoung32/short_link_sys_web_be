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
	"short_link_sys_web_be/link_gen"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/middleware"
	"short_link_sys_web_be/router"
	"syscall"
	"time"
)

func init() {
	log.GetLogger().Info("all module has init")
}

func main() {
	moduleLogger := log.GetLogger()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	if conf.GlobalConfig.GetString("mode") == "dev" {
		gin.SetMode(gin.DebugMode)
	} else if conf.GlobalConfig.GetString("mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(nil))
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
