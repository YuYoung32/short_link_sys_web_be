/**
 * Created by YuYoung on 2023/3/22
 * Description: server性能监控handler
 */

package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/log"
	"sync"
	"time"
)

var (
	staticInfo StaticInfo

	info1MinWrapper struct {
		mutex    sync.Mutex
		info1Min Info1Min
	}
	info1SWrapper struct {
		mutex  sync.Mutex
		info1s Info1s
	}

	transferGap = time.Duration(conf.GlobalConfig.GetInt64("handler.server.transferGap")) * time.Millisecond

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Init() {
	fetchInfoFromCore()
}

// fetchInfoFromCore 从转发服务器获取数据
func fetchInfoFromCore() {
	moduleLogger := log.MainLogger.WithField("func", "fetchInfoFromCore")

	wsURL := "ws://" + conf.GlobalConfig.GetString("core.host") + ":" + conf.GlobalConfig.GetString("core.port") + "/"
	moduleLogger.Info("wsURL", wsURL)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		moduleLogger.Error("dial ws failed: ", err)
		return
	}

	// auth 验证
	go func() {
		for {
			err := conn.WriteMessage(websocket.TextMessage, []byte(conf.GlobalConfig.GetString("core.auth")))
			if err != nil {
				moduleLogger.Error("write auth failed: ", err)
				return
			}
		}
	}()

	// 读取数据
	go func() {
		// 首次读取, 读取静态数据
		msgType, staticInfoBytes, err := conn.ReadMessage()
		if err != nil {
			moduleLogger.Error("read message failed: ", err)
			return
		}
		if msgType != websocket.TextMessage {
			moduleLogger.Error("auth failed")
			return
		}

		err = json.Unmarshal(staticInfoBytes, &staticInfo)
		if err != nil {
			moduleLogger.Error("unmarshal static info1MinWrapper failed: ", err)
		}
		PrintFields(staticInfo)

		// 之后读取实时数据
		go func() {
			for {
				_, dynamicInfoBytes, err := conn.ReadMessage()
				if err != nil {
					moduleLogger.Error("read dynamic message failed: ", err)
					return
				}

				info1SWrapper.mutex.Lock()
				if err = json.Unmarshal(dynamicInfoBytes, &info1SWrapper.info1s); err != nil {
					moduleLogger.Error("unmarshal dynamic info1SWrapper failed: ", err)
				}
				info1MinWrapper.mutex.Lock()
				pushAndPopArr(&info1MinWrapper.info1Min, info1SWrapper.info1s)
				info1MinWrapper.mutex.Unlock()
				info1SWrapper.mutex.Unlock()
				PrintFields(info1SWrapper.info1s)
			}
		}()
	}()
}

// pushAndPopArr 用于将实时数据推入数组，并将数组中的数据向前移动一位
func pushAndPopArr(info1Min *Info1Min, data Info1s) {
	for i := 0; i < 59; i++ {
		info1Min.CPUUsageRatioMin[i] = info1Min.CPUUsageRatioMin[i+1]
		info1Min.MemUsageMin[i] = info1Min.MemUsageMin[i+1]
		info1Min.DiskReadMin[i] = info1Min.DiskReadMin[i+1]
		info1Min.DiskWriteMin[i] = info1Min.DiskWriteMin[i+1]
		info1Min.NetRecvMin[i] = info1Min.NetRecvMin[i+1]
		info1Min.NetSendMin[i] = info1Min.NetSendMin[i+1]
	}
	info1Min.CPUUsageRatioMin[59] = data.CPUUsageRatioSec
	info1Min.MemUsageMin[59] = data.MemUsageSec
	info1Min.DiskReadMin[59] = data.DiskReadSec
	info1Min.DiskWriteMin[59] = data.DiskWriteSec
	info1Min.NetRecvMin[59] = data.NetRecvSec
	info1Min.NetSendMin[59] = data.NetSendSec
}

func Info1MinListHandler(ctx *gin.Context) {
	ctx.Set("func", "info_1_min_list_handler")
	info1MinWrapper.mutex.Lock()
	ctx.JSON(http.StatusOK, info1MinWrapper.info1Min)
	info1MinWrapper.mutex.Unlock()
}

func RealtimeDataHandler(ctx *gin.Context) {
	ModuleLogger := log.MainLogger.WithField("func", "realtime_data_handler")

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ModuleLogger.Error(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			ModuleLogger.Error(err)
		}
	}(conn)

	// 每秒发送数据
	for {
		info1SWrapper.mutex.Lock()
		err = conn.WriteJSON(info1SWrapper.info1s)
		info1SWrapper.mutex.Unlock()
		if err != nil {
			ModuleLogger.Error(err)
			return
		}
		time.Sleep(transferGap)
	}
}

// PrintFields 递归打印结构体中所有的字段和相应的值
func PrintFields(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 如果该字段是结构体，则递归打印
		if field.Kind() == reflect.Struct {
			fmt.Printf("%s:\n", fieldType.Name)
			PrintFields(field.Addr().Interface())
		} else {
			fmt.Printf("%s: %v\n", fieldType.Name, field.Interface())
		}
	}
	fmt.Println()
}
