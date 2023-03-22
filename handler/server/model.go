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

type RealtimeData struct {
	CPUUsageRatioLastSec  int `json:"cpuUsageRatioLastSec"`
	MemUsageRatioLastSec  int `json:"memUsageRatioLastSec"`
	DiskUsageRatioLastSec int `json:"diskUsageRatioLastSec"`
}

type StaticInfo struct {
	//单位: MB
	MemTotalSize  int `json:"memTotalSize"`
	DiskTotalSize int `json:"diskTotalSize"`
}

type CPURatioList [60]int
type MemRatioList [60]int
type DiskRatioList [60]int

type TTLList [60]int
