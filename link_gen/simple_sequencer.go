/**
 * Created by YuYoung on 2023/4/17
 * Description: 简单自增
 */

package link_gen

import (
	"os"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/log"
	"strconv"
)

var (
	code                               uint64 = 1000000000
	startCodeKey                              = "handler.link.simpleSeq.start"
	simpleSequencerPersistenceFilename        = "link_gen/simple_sequencer_persistence"
)

type SimpleSequencer struct{}

func simpleSequencerInit() {
	logger := log.GetLogger()
	if conf.GlobalConfig.GetString("mode") == "dev" {
		if conf.GlobalConfig.IsSet(startCodeKey) {
			code = conf.GlobalConfig.GetUint64(startCodeKey)
		} else {
			logger.Error("config file not set handler.link.simpleSeq.start, use default value 1000000000")
		}
		return
	}
	filename := simpleSequencerPersistenceFilename

	bytes, err := os.ReadFile(filename)
	if err != nil {
		logger.Errorf("Failed to read file %s: %v", filename, err)
	}
	code, err = strconv.ParseUint(string(bytes), 10, 64)
	logger.Infof("Read code %d from file %s", code, filename)
}

func simpleSequencerTerminate() {
	logger := log.GetLogger()
	filename := simpleSequencerPersistenceFilename
	err := os.WriteFile(filename, []byte(strconv.FormatUint(code, 10)), 0666)
	if err != nil {
		logger.Errorf("Failed to write file %s: %v", filename, err)
		return
	}
}

func (SimpleSequencer) GenLink(_ string) string {
	newCode := code
	code++
	return uint64ToShortLink(newCode)
}

func (SimpleSequencer) GetType() AlgorithmType {
	return SeqType
}
