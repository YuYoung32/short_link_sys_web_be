/**
 * Created by YuYoung on 2023/4/12
 * Description: 数据库等配置文件解析
 */

package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

var GlobalConfig = viper.New()

func Init() {
	GlobalConfig.SetConfigName("config")
	GlobalConfig.SetConfigType("yaml")
	GlobalConfig.AddConfigPath("./conf")

	GlobalConfig.WatchConfig()

	if err := GlobalConfig.ReadInConfig(); err != nil {
		fmt.Println(fmt.Errorf("Fatal error config file: %s \n", err))
		panic(err)
	}
}
