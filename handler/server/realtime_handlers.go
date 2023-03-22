/**
 * Created by YuYoung on 2023/3/22
 * Description: 用于websocket实时数据传输
 */

package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"short_link_sys_web_be/log"
	"time"
)

var rtd = RealtimeData{
	CPUUsageRatioLastSec:  20,
	MemUsageRatioLastSec:  30,
	DiskUsageRatioLastSec: 50,
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
		jsonStats, err := json.Marshal(rtd)
		if err != nil {
			ModuleLogger.Error(err)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, jsonStats)
		if err != nil {
			ModuleLogger.Error(err)
			return
		}

		rtd.CPUUsageRatioLastSec += 1
		rtd.MemUsageRatioLastSec += 1
		rtd.DiskUsageRatioLastSec += 1

		time.Sleep(transferGap)
	}
}
