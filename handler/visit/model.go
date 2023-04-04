/**
 * Created by YuYoung on 2023/3/22
 * Description: 返回时的数据结构
 */

package visit

import "short_link_sys_web_be/database"

type DetailsResponse = database.Details

type IPSourceResponse struct {
	Region string `json:"region"`
	Amount int    `json:"amount"`
}

type AmountTimeResponse struct {
	Amount []int `json:"amount"`
}
