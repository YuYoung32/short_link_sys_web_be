/**
 * Created by YuYoung on 2023/3/22
 * Description: 日志配置文件
 */

package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"short_link_sys_web_be/conf"
	"time"
)

var (
	// MainLogger 全局Logrus实例
	MainLogger *logrus.Logger
)

func init() {
	MainLogger = logrus.New()
	level := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}
	if conf.GlobalConfig.GetString("mode") == "dev" {
		MainLogger.SetLevel(level["debug"])
	} else {
		MainLogger.SetLevel(level[conf.GlobalConfig.GetString("log.level")])
	}

	MainLogger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     false,
	})

	logFilePath := conf.GlobalConfig.GetString("log.path") + "/" + "log.log"
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	MainLogger.SetOutput(io.MultiWriter(file, os.Stdout))

	GetLogger().Info("Logrus init success")
}

// Middleware 日志中间件
func Middleware(ctx *gin.Context) {
	// 开始时间
	startTime := time.Now()

	// 处理请求
	ctx.Next()

	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)

	// 请求方式
	reqMethod := ctx.Request.Method

	// 请求路由
	reqUri := ctx.Request.RequestURI

	// 状态码
	statusCode := ctx.Writer.Status()

	// 请求IP
	clientIP := ctx.ClientIP()

	//MainLogger.Infof("| %3d | %13v | %15s | %s | %s |",
	//	statusCode,
	//	latencyTime,
	//	clientIP,
	//	reqMethod,
	//	reqUri,
	//)
	MainLogger.Debugf("| %3d | %13v | %15s | %s | %s |",
		statusCode,
		latencyTime,
		clientIP,
		reqMethod,
		reqUri,
	)
}

// GetLogger 获取日志实例, WithField为获得调用方的函数名
func GetLogger() *logrus.Entry {
	// 获取调用栈信息
	pc, _, _, _ := runtime.Caller(1)
	// 获取函数名
	funcName := runtime.FuncForPC(pc).Name()
	return MainLogger.WithField("func", funcName)
}

// GetLoggerWithSkip 获取日志实例 skip=1 为调用GetLogger的函数, skip=2 为调用GetLogger的函数的上一级函数, 以此类推
func GetLoggerWithSkip(skip int) *logrus.Entry {
	// 获取调用栈信息
	pc, _, _, _ := runtime.Caller(skip)
	// 获取函数名
	funcName := runtime.FuncForPC(pc).Name()
	return MainLogger.WithField("func", funcName)
}
