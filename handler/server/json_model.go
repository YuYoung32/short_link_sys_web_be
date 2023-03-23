/**
 * Created by YuYoung on 2023/3/22
 * Description: 整合每个接口的JSON数据结构
 */

package server

//type RealTimeData struct {
//	ShortLink  string `json:"shortLink"`
//	LongLink   string `json:"longLink"`
//	CreateTime string `json:"createTime"`
//}

type Info1s struct {
	CPUUsageRatioSec  int `json:"cpuUsageRatioLastSec"`
	MemUsageRatioSec  int `json:"memUsageRatioLastSec"`
	DiskUsageRatioSec int `json:"diskUsageRatioLastSec"`
}

type Info1Min struct {
	CPUUsageRatioSec  [60]int `json:"cpuUsageRatioLastMin"`
	MemUsageRatioSec  [60]int `json:"memUsageRatioLastMin"`
	DiskUsageRatioSec [60]int `json:"diskUsageRatioLastMin"`
}

type InfoXhr struct {
	TimePoints        [60]int `json:"xHourTimePoints"`
	CPUUsageRatioMin  [60]int `json:"cpuUsageRatioLastXHour"`
	MemUsageRatioMin  [60]int `json:"memUsageRatioLastXHour"`
	DiskUsageRatioMin [60]int `json:"diskUsageRatioLastXHour"`
	TTLMin            [60]int `json:"ttlLastXHour"`
}

type StaticInfo struct {
	//单位: MB
	MemTotalSize  int `json:"memTotalSize"`
	DiskTotalSize int `json:"diskTotalSize"`
}
