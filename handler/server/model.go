/**
 * Created by YuYoung on 2023/3/22
 * Description: 整合每个接口的JSON数据结构
 */

package server

type Info1s struct {
	RunningTime      int `json:"runningTime"`
	CPUUsageRatioSec int `json:"cpuUsageRatioLastSec"`
	CPUFreqSec       int `json:"cpuFreqLastSec"`

	MemUsageSec int `json:"memUsageLastSec"`
	MemAvailSec int `json:"memAvailLastSec"`
	SwapUsage   int `json:"swapUsageLastSec"`

	DiskReadSec  int `json:"diskReadLastSec"`
	DiskWriteSec int `json:"diskWriteLastSec"`
	DiskUsageSec int `json:"diskUsageLastSec"`
	DiskAvailSec int `json:"diskAvailLastSec"`

	NetRecvSec int `json:"netRecvLastSec"`
	NetSendSec int `json:"netSendLastSec"`

	TTLSec int `json:"ttlLastSec"`
}

type Info1Min struct {
	CPUUsageRatioMin [60]int `json:"cpuUsageRatioLastMin"`
	MemUsageMin      [60]int `json:"memUsageLastMin"`
	DiskReadMin      [60]int `json:"diskReadLastMin"`
	DiskWriteMin     [60]int `json:"diskWriteLastMin"`
	NetRecvMin       [60]int `json:"netRecvLastMin"`
	NetSendMin       [60]int `json:"netSendLastMin"`
}

type CPUStaticInfo struct {
	Name      string `json:"name"`
	CoreNum   int    `json:"coreNum"`
	ThreadNum int    `json:"threadNum"`
	CacheSize int    `json:"cacheSize"`
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
