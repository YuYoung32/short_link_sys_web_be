/**
 * Created by YuYoung on 2023/3/22
 * Description: 整合每个接口的JSON数据结构
 */

package visit

type Details struct {
	LongUrl   string `json:"longUrl"`
	ShortUrl  string `json:"shortUrl"`
	IP        string `json:"IP"`
	Region    string `json:"region"`
	OS        string `json:"OS"`
	Timestamp string `json:"timestamp"`
}

type IPSource struct {
	Region string `json:"region"`
	Amount int    `json:"amount"`
}

type AmountTime struct {
	Amount []int `json:"amount"`
}
