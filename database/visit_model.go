/**
 * Created by YuYoung on 2023/4/1
 * Description: 访问 ORM Model
 */

package database

// LinkVisit 访问记录, 此为视图, 需要手动执行sql创建
type LinkVisit struct {
	ShortLink string `json:"shortLink"`
	LongLink  string `json:"longLink"`
	Comment   string `json:"comment"`
	IP        string `json:"ip"`
	Region    string `json:"region"`
	VisitTime int64  `json:"visitTime"`
}
