/**
 * Created by YuYoung on 2023/4/4
 * Description: 时间处理相关的工具函数
 */

package utils

import (
	"errors"
	"strconv"
	"time"
)

// GetTodayBeginTime 获取今天的开始时间戳
func GetTodayBeginTime() int64 {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location()).Unix()
}

// ConvertAndCheckTimeGroup 将时间字符串转换为时间戳并检查是否合法(beginStr<=endStr && endStr<=当前时间)
func ConvertAndCheckTimeGroup(beginStr string, endStr string) (begin int64, end int64, err error) {
	if beginStr == "" || endStr == "" {
		err = errors.New("miss args")
		return
	}
	if begin, err = strconv.ParseInt(beginStr, 10, 64); err != nil {
		return
	}
	if end, err = strconv.ParseInt(endStr, 10, 64); err != nil {
		return
	}
	if begin > end || end > time.Now().Unix() {
		err = errors.New("invalid args")
		return
	}
	return
}
