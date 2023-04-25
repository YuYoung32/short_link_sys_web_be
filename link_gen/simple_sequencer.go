/**
 * Created by YuYoung on 2023/4/17
 * Description: 简单自增
 */

package link_gen

import (
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/log"
)

var (
	code         uint64 = 1000000000
	startCodeKey        = "handler.link.simpleSeq.start"
)

type SimpleSequencer struct{}

func simpleSequencerInit() {
	if conf.GlobalConfig.IsSet(startCodeKey) {
		code = conf.GlobalConfig.GetUint64(startCodeKey)
	} else {
		log.GetLogger().Error("config file not set handler.link.simpleSeq.start, use default value 1000000000")
	}
}

func simpleSequencerTerminate() {
	conf.GlobalConfig.Set(startCodeKey, code)
	err := conf.GlobalConfig.WriteConfig()
	if err != nil {
		log.GetLogger().Errorf("write %v failed, value %v", startCodeKey, code)
		return
	}
}

func (SimpleSequencer) GenLink(_ string) string {
	mutex.Lock()
	newCode := code
	code++
	mutex.Unlock()
	return uint64ToShortLink(newCode)
}

func (SimpleSequencer) GetType() AlgorithmType {
	return SeqType
}
