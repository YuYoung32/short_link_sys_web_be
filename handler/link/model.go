/**
 * Created by YuYoung on 2023/3/22
 * Description: 整合每个接口的JSON数据结构
 */

package link

import "short_link_sys_web_be/database"

type AmountTotal struct {
	Amount int64 `json:"amountTotal"`
}

type DetailsListResponse struct {
	Links      []database.Link `json:"links"`
	LinksTotal int64           `json:"linksTotal"`
}
