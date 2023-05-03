/**
 * Created by YuYoung on 2023/5/2
 * Description: 公用的函数
 */

package perf

import (
	"fmt"
	"short_link_sys_web_be/link_gen"
	"strconv"
	"time"
)

func GenLinks(amount int, inter link_gen.LinkGen) []string {
	links := make([]string, amount)
	for i := 0; i < amount; i++ {
		links[i] = inter.GenLink("www.baidu.com" + strconv.Itoa(i))
	}
	return links
}

func IntsToFloat64(d []int) []float64 {
	res := make([]float64, len(d))
	for i := 0; i < len(d); i++ {
		res[i] = float64(d[i])
	}
	return res
}

func DurationsToFloat64(d []time.Duration) []float64 {
	res := make([]float64, len(d))
	for i := 0; i < len(d); i++ {
		res[i] = float64(d[i].Nanoseconds()) / float64(time.Millisecond)
	}
	return res
}

func PrintCommaSeqFloatSlice(s []float64) string {
	var str string
	for i := 0; i < len(s); i++ {
		if i == len(s)-1 {
			str += fmt.Sprintf("%f", s[i])
			return str
		}
		str += fmt.Sprintf("%f, ", s[i])
	}
	return str
}
