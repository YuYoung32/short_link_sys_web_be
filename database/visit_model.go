/**
 * Created by YuYoung on 2023/4/1
 * Description: 访问 ORM Model
 */

package database

type Visit struct {
	ShortLink string `json:"shortLink" gorm:"primaryKey"`
	LongLink  string `json:"longLink"`
	Comment   string `json:"comment"`
	IP        string `json:"ip"`
	Region    string `json:"region"`
	VisitTime int64  `json:"visitTime" gorm:"autoCreateTime"`
}
