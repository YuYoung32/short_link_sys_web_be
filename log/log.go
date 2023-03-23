/**
 * Created by YuYoung on 2023/3/22
 * Description: 日志配置文件
 */

package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	// MainLogger 全局Logrus实例
	MainLogger = logrus.New()
)

type TempLogConf struct {
	Level   string `json:"level"`
	DicPath string `json:"file_path"`
}

func GetLogConf() TempLogConf {
	return TempLogConf{
		Level:   "debug",
		DicPath: "log/",
	}
}

// 配置Logrus
func init() {
	logConf := GetLogConf()
	level := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}
	MainLogger.SetLevel(level[logConf.Level])

	MainLogger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	//logFilePath := logConf.DicPath + time.Now().Format("2006-01-02-15-04-05") + "_log.log"
	//file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//if err != nil {
	//	panic(err)
	//}

	//MainLogger.SetOutput(io.MultiWriter(file, os.Stdout))
	MainLogger.SetOutput(os.Stdout)

	MainLogger.Info("Logrus init success")
}

// Middleware 日志中间件
func Middleware(c *gin.Context) {
	// 开始时间
	startTime := time.Now()

	// 处理请求
	c.Next()

	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)

	// 请求方式
	reqMethod := c.Request.Method

	// 请求路由
	reqUri := c.Request.RequestURI

	// 状态码
	statusCode := c.Writer.Status()

	// 请求IP
	clientIP := c.ClientIP()

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