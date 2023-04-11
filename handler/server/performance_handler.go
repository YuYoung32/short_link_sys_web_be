/**
 * Created by YuYoung on 2023/3/22
 * Description: server性能监控handler
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

type info1MinMutex struct {
	mutex    sync.Mutex
	info1Min Info1Min
}

var (
	info info1MinMutex
	rtd  info1sMutex

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	realTimeDataTransfer = make(chan Info1s, 1)
)

func init() {
	go func() {
		for data := range realTimeDataTransfer {
			info.mutex.Lock()
			pushAndPopArr(&info.info1Min, data)
			info.mutex.Unlock()
		}
	}()
}

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
	ctx.Set("module", "info_1_min_list_handler")
	ctx.JSON(http.StatusOK, info.info1Min)
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

	// 产生数据并发送, 1s一次
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

		rtd.info1s.MemUsageSec = rand.Intn(1000)*1024*1024 + 5*1024*1024*1024
		rtd.info1s.MemAvailSec = testStaticInfoData.MemStaticInfo.PhysicalTotalSize - rtd.info1s.MemUsageSec
		rtd.info1s.SwapUsage = rand.Intn(1000)*1024*1024 + 1000

		rtd.info1s.DiskUsageSec = rand.Intn(101)
		rtd.info1s.CPUFreqSec = rand.Intn(1000) + 1000
		rtd.info1s.RunningTime = rtd.info1s.RunningTime + 1

		rtd.info1s.DiskReadSec = rand.Intn(1000) + 1000
		rtd.info1s.DiskWriteSec = rand.Intn(1000) + 1000
		rtd.info1s.DiskUsageSec = rand.Intn(1000)*1024*1024 + 20*1024*1024*1024
		rtd.info1s.DiskAvailSec = testStaticInfoData.DiskStaticInfo.DiskTotalSize - rtd.info1s.DiskUsageSec

		rtd.info1s.NetRecvSec = rand.Intn(1000) + 1000
		rtd.info1s.NetSendSec = rand.Intn(1000) + 1000

		rtd.info1s.TTLSec = rand.Intn(2000)

		rtd.mutex.Unlock()

		realTimeDataTransfer <- rtd.info1s
		time.Sleep(transferGap)
	}
}
