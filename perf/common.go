/**
 * Created by YuYoung on 2023/5/2
 * Description: 公用的函数
 */

package perf

import (
	"bufio"
	"fmt"
	"os"
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

// DurationsToFloat64 将time.Duration转换为float64 毫秒, 保留精度
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

func AvgFloat64(s []float64) float64 {
	var sum float64
	for i := 0; i < len(s); i++ {
		sum += s[i]
	}
	return sum / float64(len(s))
}

func SaveFloat64SliceToFile(arr []float64, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(file)

	writer := bufio.NewWriter(file)
	for _, f := range arr {
		_, err := fmt.Fprintf(writer, "%f\n", f)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func ReadFileToFloat64Slice(path string) ([]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(file)

	scanner := bufio.NewScanner(file)

	var arr []float64
	for scanner.Scan() {
		line := scanner.Text()
		f, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, err
		}
		arr = append(arr, f)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return arr, nil
}
