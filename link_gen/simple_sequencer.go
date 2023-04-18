/**
 * Created by YuYoung on 2023/4/17
 * Description: 简单自增
 */

package link_gen

var code uint64 = 1000000000

type SimpleSequencer struct{}

func (SimpleSequencer) GenLink(s string) string {
	mutex.Lock()
	code++
	newCode := code
	mutex.Unlock()
	return uint64ToShortLink(newCode)
}

func (SimpleSequencer) GetType() AlgorithmType {
	return SeqType
}
