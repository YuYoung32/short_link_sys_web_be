/**
 * Created by YuYoung on 2023/4/17
 * Description: FNV hash算法
 */

package link_gen

import (
	"hash/fnv"
)

type fnvHash struct{}

var fnvHash64 = fnv.New64()

func (fnvHash) GenLink(s string) string {
	fnvHash64.Reset()
	_, err := fnvHash64.Write([]byte(s))
	if err != nil {
		return err.Error()
	}
	return uint64ToShortLink(fnvHash64.Sum64())
}
