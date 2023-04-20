/**
 * Created by YuYoung on 2023/3/22
 * Description: 返回时的数据结构
 */

package visit

import "short_link_sys_web_be/database"

type IPSourceResponse struct {
	Region []string `json:"region"`
	Amount []int    `json:"amount"`
}

type DetailsListResponse struct {
	VisitDetails       []database.LinkVisit `json:"visitDetails"`
	VisitDetailsAmount int64                `json:"visitDetailsAmount"`
}

type AmountResponse struct {
	Amount int `json:"amount"`
}

type StaticsListResponse struct {
	VisitAmountList []int `json:"visitAmount"`
	IPAmountList    []int `json:"ipAmount"`
}
