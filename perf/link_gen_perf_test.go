/**
 * Created by YuYoung on 2023/5/2
 * Description: 链接生成速率测试
 */

package perf

import (
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/wcharczuk/go-chart"
	"gonum.org/v1/gonum/stat"
	"os"
	"short_link_sys_web_be/link_gen"
	"strconv"
	"testing"
	"time"
)

var amountSeq []int
var amount = 100
var gap = 1000

func TestMain(m *testing.M) {
	for i := 1; i <= amount; i++ {
		amountSeq = append(amountSeq, i*gap)
	}
	m.Run()
}

// getTimeRecord 获取生成链接的时间记录
func getTimeRecord(inter link_gen.LinkGen) []time.Duration {
	var timeRecorder = make([]time.Duration, len(amountSeq))
	for i := 0; i < len(amountSeq); i++ {
		start := time.Now()
		GenLinks(amountSeq[i], inter)
		timeRecorder[i] = time.Since(start)
	}
	return timeRecorder
}

// TestLinkGen 测试链接生成速率
func TestLinkGenSpeed(t *testing.T) {
	ssRec := getTimeRecord(link_gen.SimpleSequencer{})
	sfsRec := getTimeRecord(link_gen.SnowflakeSequencer{})
	fnvhRec := getTimeRecord(link_gen.FNVHash{})
	mmhRec := getTimeRecord(link_gen.MurmurHash{})
	xxhRec := getTimeRecord(link_gen.XXHash{})

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Number of Links",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				i, _ := v.(float64)
				return strconv.Itoa(int(i))
			},
		},
		YAxis: chart.YAxis{
			Name:      "Consumption(ms)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},

		Series: []chart.Series{
			chart.ContinuousSeries{
				Name: "Simple Sequencer",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: DurationsToFloat64(ssRec),
			},
			chart.ContinuousSeries{
				Name: "Snowflake Sequencer",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(1),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: DurationsToFloat64(sfsRec),
			},
			chart.ContinuousSeries{
				Name: "FNV Hash",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(2),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: DurationsToFloat64(fnvhRec),
			},
			chart.ContinuousSeries{
				Name: "Murmur Hash",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(3),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: DurationsToFloat64(mmhRec),
			},
			chart.ContinuousSeries{
				Name: "XX Hash",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(4),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: DurationsToFloat64(xxhRec),
			},
		},
	}
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	f, err := os.Create("perf/link_gen_res/link_gen_speed.png")
	if err != nil {
		t.Error(err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}
	}(f)
	err = graph.Render(chart.PNG, f)
	if err != nil {
		t.Error(err)
		return
	}

	a, b := stat.LinearRegression(IntsToFloat64(amountSeq), DurationsToFloat64(ssRec), nil, true)
	t.Logf("simple sequencer y=%f+%fx", a, b)

	a, b = stat.LinearRegression(IntsToFloat64(amountSeq), DurationsToFloat64(sfsRec), nil, true)
	t.Logf("snowflake sequencer y=%f+%fx", a, b)

	a, b = stat.LinearRegression(IntsToFloat64(amountSeq), DurationsToFloat64(fnvhRec), nil, true)
	t.Logf("fnv hash y=%f+%fx", a, b)

	a, b = stat.LinearRegression(IntsToFloat64(amountSeq), DurationsToFloat64(mmhRec), nil, true)
	t.Logf("murmur hash y=%f+%fx", a, b)

	a, b = stat.LinearRegression(IntsToFloat64(amountSeq), DurationsToFloat64(xxhRec), nil, true)
	t.Logf("xx hash y=%f+%fx", a, b)
}

// TestHashCollision 测试Hash算法生成的哈希碰撞
func TestHashCollision(t *testing.T) {
	getHashCollisionRecord := func(inter link_gen.LinkGen) []float64 {
		var hashRecorder = make([]float64, len(amountSeq))
		for i := 0; i < len(amountSeq); i++ {
			t.Log("running on", amountSeq[i], "links")
			shortLinkBF := bloom.NewWithEstimates(1000000, 0.01)
			cnt := 0
			tmpSLList := make([]string, amountSeq[i])
			for j := 0; j < amountSeq[i]; j++ {
				shortLink := inter.GenLink("https://baidu.com" + strconv.Itoa(j))
				tmpSLList[j] = shortLink
				if shortLinkBF.TestAndAddString(shortLink) {
					// 排除假阳性
					for k := 0; k < j; k++ {
						if tmpSLList[k] == shortLink {
							cnt++
							break
						}
					}
				}
			}
			hashRecorder[i] = float64(cnt) / float64(amountSeq[i])
		}
		return hashRecorder
	}
	fnvhRec := getHashCollisionRecord(link_gen.FNVHash{})
	mmhRec := getHashCollisionRecord(link_gen.MurmurHash{})
	xxhRec := getHashCollisionRecord(link_gen.XXHash{})

	t.Log("fnv Hash", fnvhRec)
	t.Log("murmur Hash", mmhRec)
	t.Log("xx Hash", xxhRec)
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Number of Links",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				i, _ := v.(float64)
				return fmt.Sprintf("%.2f", i)
			},
		},
		YAxis: chart.YAxis{
			Name:      "Consumption(ms)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},

		Series: []chart.Series{
			chart.ContinuousSeries{
				Name: "FNV Hash",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: fnvhRec,
			},
			chart.ContinuousSeries{
				Name: "Murmur Hash",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(1),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: mmhRec,
			},
			chart.ContinuousSeries{
				Name: "XX Hash",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(2),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: xxhRec,
			},
		},
	}
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	f, err := os.Create("perf/link_gen_res/link_gen_hash_collision.png")
	if err != nil {
		t.Error(err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}
	}(f)
	err = graph.Render(chart.PNG, f)
	if err != nil {
		t.Error(err)
		return
	}
}

//func GenHashLinks(amount int, inter LinkGen) []string {
//	//var shortLinkBF = bloom.NewWithEstimates(
//	//	1000000,
//	//	0.01)
//	//genUniqueLink := func(longLink string) string {
//	//	var shortLink string
//	//	for {
//	//		shortLink = inter.GenLink(longLink)
//	//		if !shortLinkBF.TestString(shortLink) {
//	//			break
//	//		}
//	//		longLink = longLink + strconv.Itoa(rand.Int())
//	//		shortLink = inter.GenLink(longLink)
//	//	}
//	//	return shortLink
//	//}
//
//	links := make([]string, amount)
//	for i := 0; i < amount; i++ {
//		links[i] = inter.GenLink("www.baidu.com" + strconv.Itoa(i))
//	}
//	return links
//}
