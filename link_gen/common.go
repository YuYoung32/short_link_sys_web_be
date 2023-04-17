/**
 * Created by YuYoung on 2023/4/17
 * Description: 短链生成函数规范
 */

package link_gen

func Init() {
	SnowflakeInit()
}

type LinkGen interface {
	GenLink(string) string
}

const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func uint64ToBase62(n uint64) string {
	var result string
	for n > 0 {
		result = string(base62[n%62]) + result
		n = n / 62
	}
	return result
}

func fillZero(str string) string {
	for len(str) < 6 {
		str = "0" + str
	}
	return str
}

func uint64ToShortLink(n uint64) string {
	return fillZero(uint64ToBase62(n))
}
