/**
 * Created by YuYoung on 2023/3/22
 * Description: 用于websocket实时数据传输
 */

package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"short_link_sys_web_be/log"
	"sync"
	"time"
)

type info1sMutex struct {
	mutex  sync.Mutex
	info1s Info1s
}

var (
	rtd info1sMutex

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	realTimeDataTransfer = make(
		chan Info1s, 1,
	)
)

func RealtimeDataHandler(ctx *gin.Context) {
	ModuleLogger := log.MainLogger.WithField("module", "realtime_data_handler")
	transferGap := 1 * time.Second

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

	for {
		rand.Seed(time.Now().UnixNano())

		jsonStats, err := json.Marshal(rtd.info1s)
		if err != nil {
			ModuleLogger.Error(err)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, jsonStats)
		if err != nil {
			ModuleLogger.Error(err)
			return
		}

		rtd.mutex.Lock()
		// 竞争：远端读数据，本地读取再发送
		rtd.info1s.CPUUsageRatioSec = rand.Intn(101)
		rtd.info1s.MemUsageRatioSec = rand.Intn(101)
		rtd.info1s.DiskUsageRatioSec = rand.Intn(101)
		rtd.mutex.Unlock()

		realTimeDataTransfer <- rtd.info1s
		time.Sleep(transferGap)
	}
}
