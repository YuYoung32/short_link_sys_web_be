/**
 * Created by YuYoung on 2023/4/17
 * Description: 简单自增
 */

package link_gen

var start = uint64(1000000000000000000)

type SimpleSequencer struct{}

func (SimpleSequencer) GenLink(s string) string {
	start++
	return uint64ToShortLink(start)
}
