/**
 * Created by YuYoung on 2023/4/17
 * Description: 短链生成函数规范
 */

package link_gen

import (
	"short_link_sys_web_be/conf"
	"strings"
	"sync"
)

var (
	minLength int
	mutex     = sync.Mutex{} // 互斥锁, 保证并发安全
)

type AlgorithmType int

const (
	HashType AlgorithmType = iota
	SeqType
)

func Init() {
	minLength = conf.GlobalConfig.GetInt("handler.link.minLength")
	SnowflakeInit()
	simpleSequencerInit()
}

func Terminate() {
	simpleSequencerTerminate()
}

type LinkGen interface {
	GenLink(string) string
	GetType() AlgorithmType
}

const base = "hibnP8XAcde7qrsIFzMUaZgHVJ3f0STu169WjklmGy4BCDLEQRvwtY25xKopNO"
const baseLen = uint64(len(base))

func uint64ToBase(n uint64) string {
	var result string
	for n > 0 {
		result = string(base[n%baseLen]) + result
		n = n / baseLen
	}
	return result
}

func baseToUint64(s string) uint64 {
	var result uint64
	for _, c := range s {
		result = result*baseLen + uint64(strings.IndexByte(base, byte(c)))
	}
	return result
}

func fillZero(str string) string {
	for len(str) < minLength {
		str = "0" + str
	}
	return str
}

func uint64ToShortLink(n uint64) string {
	return fillZero(uint64ToBase(n))
}
