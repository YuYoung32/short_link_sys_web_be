/**
 * Created by YuYoung on 2023/4/17
 * Description: FNV hash算法
 */

package link_gen

import (
	"hash/fnv"
)

type FNVHash struct{}

var FNVHash64 = fnv.New64()

func (FNVHash) GenLink(s string) string {
	FNVHash64.Reset()
	_, err := FNVHash64.Write([]byte(s))
	if err != nil {
		return err.Error()
	}
	return uint64ToShortLink(FNVHash64.Sum64())
}
