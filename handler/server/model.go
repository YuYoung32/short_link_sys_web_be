/**
 * Created by YuYoung on 2023/3/22
 * Description: 整合每个接口的JSON数据结构
 */

package server

type Info1s struct {
	CPUUsageRatioSec int   `json:"cpuUsageRatioLastSec"`
	CPURunningTime   int64 `json:"cpuRunningTime"`

	MemUsageSec uint64 `json:"memUsageLastSec"`
	MemAvailSec uint64 `json:"memAvailLastSec"`
	SwapUsage   uint64 `json:"swapUsageLastSec"`

	DiskReadSec  uint64 `json:"diskReadLastSec"`
	DiskWriteSec uint64 `json:"diskWriteLastSec"`
	DiskUsageSec uint64 `json:"diskUsageLastSec"`
	DiskAvailSec uint64 `json:"diskAvailLastSec"`

	NetRecvSec uint64 `json:"netRecvLastSec"`
	NetSendSec uint64 `json:"netSendLastSec"`

	//TTLSec int `json:"ttlLastSec"`
}

type Info1Min struct {
	CPUUsageRatioMin [60]int    `json:"cpuUsageRatioLastMin"`
	MemUsageMin      [60]uint64 `json:"memUsageLastMin"`
	DiskReadMin      [60]uint64 `json:"diskReadLastMin"`
	DiskWriteMin     [60]uint64 `json:"diskWriteLastMin"`
	NetRecvMin       [60]uint64 `json:"netRecvLastMin"`
	NetSendMin       [60]uint64 `json:"netSendLastMin"`
}

type CPUStaticInfo struct {
	Name      string `json:"name"`
	CoreNum   int    `json:"coreNum"`
	ThreadNum int    `json:"threadNum"`
	CacheSize int    `json:"cacheSize"`
	Speed     int    `json:"speed"`
	StartTime int64  `json:"startTime"`
}

type MemStaticInfo struct {
	PhysicalTotalSize int `json:"physicalTotalSize"`
	SwapTotalSize     int `json:"swapTotalSize"`
}

type DiskStaticInfo struct {
	DiskTotalSize int `json:"diskTotalSize"`
}

type NetStaticInfo struct {
	IPv4 string `json:"ipv4"`
	MAC  string `json:"mac"`
}

type StaticInfo struct {
	CPUStaticInfo  CPUStaticInfo  `json:"cpuStaticInfo"`
	MemStaticInfo  MemStaticInfo  `json:"memStaticInfo"`
	DiskStaticInfo DiskStaticInfo `json:"diskStaticInfo"`
	NetStaticInfo  NetStaticInfo  `json:"netStaticInfo"`
}
