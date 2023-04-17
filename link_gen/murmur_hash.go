/**
 * Created by YuYoung on 2023/4/17
 * Description: murmurHash算法
 */

package link_gen

import "github.com/spaolacci/murmur3"

type MurmurHash struct{}

func (m MurmurHash) GenLink(s string) string {
	murmur3.New64()
	hash := murmur3.Sum64([]byte(s))
	return uint64ToShortLink(hash)
}
