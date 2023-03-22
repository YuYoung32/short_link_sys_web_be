/**
 * Created by YuYoung on 2023/3/22
 * Description: 整合每个接口的JSON数据结构
 */

package link

type Details struct {
	ShortLink  string `json:"shortLink"`
	LongLink   string `json:"longLink"`
	CreateTime string `json:"createTime"`
}

type TotalAmount int
