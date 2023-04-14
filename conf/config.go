/**
 * Created by YuYoung on 2023/4/12
 * Description: 数据库等配置文件解析
 */

package conf

import (
	"github.com/spf13/viper"
)

var GlobalConfig = viper.New()

func Init() {
	GlobalConfig.SetConfigName("config")
	GlobalConfig.SetConfigType("yaml")
	GlobalConfig.AddConfigPath("C:\\Users\\29011\\GolandProjects\\short_link_sys_web_be\\conf\\")
	if err := GlobalConfig.ReadInConfig(); err != nil {
		panic(err)
	}
}
