/**
 * Created by YuYoung on 2023/4/17
 * Description: xxHash算法
 */

package link_gen

import "github.com/cespare/xxhash"

type xxHash struct{}

func (xxHash) GenLink(s string) string {
	hash := xxhash.Sum64String(s)
	return uint64ToShortLink(hash)
}
