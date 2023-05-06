/**
 * Created by YuYoung on 2023/5/4
 * Description: 转发服务测试
 */

package perf

import (
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"gonum.org/v1/gonum/stat"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func genSeq(n int) []float64 {
	var res []float64
	for i := 1; i <= n; i++ {
		res = append(res, float64(i))
	}
	return res
}

func saveForwardGraph1(x, y []float64, yName, path string) error {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Number of Concurrency",
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
			ValueFormatter: func(v interface{}) string {
				i, _ := v.(float64)
				return strconv.Itoa(int(i))
			},
		},

		Series: []chart.Series{
			chart.ContinuousSeries{
				Name: yName,
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0),
				},
				XValues: x,
				YValues: y,
			},
		},
	}
	//graph.Elements = []chart.Renderable{
	//	chart.Legend(&graph),
	//}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	err = graph.Render(chart.PNG, f)
	if err != nil {
		return err
	}
	return nil
}

func saveForwardGraph3(x, y1, y2, y3 []float64, y1Name, y2Name, y3Name, path string) error {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Number of Concurrency",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				i, _ := v.(float64)
				return strconv.Itoa(int(i))
			},
		},
		YAxis: chart.YAxis{
			Name:      "Average Consumption(ms)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				i, _ := v.(float64)
				return strconv.Itoa(int(i))
			},
		},

		Series: []chart.Series{
			chart.ContinuousSeries{
				Name: y1Name,
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0),
				},
				XValues: x,
				YValues: y1,
			},
			chart.ContinuousSeries{
				Name: y2Name,
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(1),
				},
				XValues: x,
				YValues: y2,
			},
			chart.ContinuousSeries{
				Name: y3Name,
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(2),
				},
				XValues: x,
				YValues: y3,
			},
		},
	}
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	err = graph.Render(chart.PNG, f)
	if err != nil {
		return err
	}
	return nil
}

// TestRightURLConcurrencyOnce 测试单次并发
func TestRightURLConcurrencyOnce(t *testing.T) {
	url := "http://rdtws.me/15hyTK"
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	concurrencyNum := 100 // 并发数

	var wg sync.WaitGroup
	var record = make([]time.Duration, concurrencyNum)
	wg.Add(concurrencyNum)
	for i := 0; i < concurrencyNum; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
			}()

			start := time.Now()
			_, err := client.Get(url)
			record[i] = time.Since(start)
			if err != nil {
				t.Log(err)
			}
		}(i)
	}

	// 等待所有请求完成
	wg.Wait()

	//err := saveForwardGraph1(genSeq(concurrencyNum), DurationsToFloat64(record), "Right", "perf/forward_res/forward_url_concurrency_once.png")
	err := saveScatterGraph(genSeq(concurrencyNum), DurationsToFloat64(record), "perf/forward_res/forward_url_concurrency_once_scatter.png")
	if err != nil {
		t.Error(err)
		return
	}
}

// TestRightURLConcurrency 测试多次请求, 次数递增取均值
func TestRightURLConcurrency(t *testing.T) {
	gap := 20
	amount := 50
	var ccSeq []int
	for i := 0; i < amount; i++ {
		ccSeq = append(ccSeq, (i+1)*gap)
	}

	//url := "http://rdtws.me/15hyTK"
	url := "http://localhost:8090/i8y540"
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var record = make([]float64, amount)
	for k := 0; k < amount; k++ {
		concurrencyNum := ccSeq[k] // 并发数
		t.Log("concurrencyNum: ", concurrencyNum)
		var wg sync.WaitGroup
		var curRecord = make([]time.Duration, concurrencyNum)
		wg.Add(concurrencyNum)
		for i := 0; i < concurrencyNum; i++ {
			go func(i int) {
				defer func() {
					wg.Done()
				}()

				start := time.Now()
				_, err := client.Get(url)
				if err != nil {
					t.Log(err)
				}
				curRecord[i] = time.Since(start)
			}(i)
		}

		// 等待所有请求完成
		wg.Wait()
		record[k] = AvgFloat64(DurationsToFloat64(curRecord))
		time.Sleep(time.Second * 1)
	}
	err := SaveFloat64SliceToFile(record, "perf/forward_res/forward_url_concurrency.txt")
	if err != nil {
		t.Log(err)
		return
	}
	err = saveForwardGraph1(IntsToFloat64(ccSeq), record, "Right", "perf/forward_res/forward_url_concurrency.png")
	if err != nil {
		t.Log(err)
		return
	}
}

func saveScatterGraph(x, y []float64, path string) error {
	viridisByY := func(xr, yr chart.Range, index int, x, y float64) drawing.Color {
		return chart.Viridis(y, yr.GetMin(), yr.GetMax())
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Number of Concurrency",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				i, _ := v.(float64)
				if int(i)%100 == 0 {
					return strconv.Itoa(int(i))
				}
				return ""
			},
		},
		YAxis: chart.YAxis{
			Name:      "Consumption(ms)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				i, _ := v.(float64)
				return strconv.Itoa(int(i))
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:             true,
					StrokeWidth:      chart.Disabled,
					DotWidth:         1,
					DotColorProvider: viridisByY,
				},
				XValues: x,
				YValues: y,
			},
		},
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	err = graph.Render(chart.PNG, f)
	if err != nil {
		return err
	}
	return nil
}

func TestScatterGraph(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	y := []float64{13, 24, 32, 41, 13, 34, 56, 23, 23, 10}
	err := saveScatterGraph(x, y, "perf/forward_res/scatter.png")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSaveDataFile(t *testing.T) {
	var seq []float64
	for i := 0; i < 100; i++ {
		seq = append(seq, float64(i))
	}
	err := SaveFloat64SliceToFile(seq, "perf/forward_res/seq.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestReadDataFile(t *testing.T) {
	seq, err := ReadFileToFloat64Slice("perf/forward_res/forward_url_concurrency_1000.txt")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(seq)
}

func TestMerge3Graph(t *testing.T) {
	var seq []int
	amount := 50
	gap := 20
	for i := 0; i < amount; i++ {
		seq = append(seq, (i+1)*gap)
	}
	s1, _ := ReadFileToFloat64Slice("perf/forward_res/forward_url_concurrency_direct.txt")
	s2, _ := ReadFileToFloat64Slice("perf/forward_res/forward_url_concurrency_redis.txt")
	s3, _ := ReadFileToFloat64Slice("perf/forward_res/forward_url_concurrency_mysql.txt")
	err := saveForwardGraph3(IntsToFloat64(seq), s1, s2, s3,
		"Test1-Direct", "Test2-Redis", "Test3-MySQL",
		"perf/forward_res/forward_url_merge.png")
	if err != nil {
		t.Error(err)
		return
	}
	s1a, s1b := stat.LinearRegression(IntsToFloat64(seq), s1, nil, false)
	s2a, s2b := stat.LinearRegression(IntsToFloat64(seq), s2, nil, false)
	s3a, s3b := stat.LinearRegression(IntsToFloat64(seq), s3, nil, false)
	t.Log("direct mse and avg:", MSE(s1a, s1b, IntsToFloat64(seq), s1), stat.Mean(s1, nil))
	t.Log("redis mse and avg:", MSE(s2a, s2b, IntsToFloat64(seq), s2), stat.Mean(s2, nil))
	t.Log("mysql mse and avg:", MSE(s3a, s3b, IntsToFloat64(seq), s3), stat.Mean(s3, nil))
	/**
	  forward_server_test.go:357: direct mse and avg: 30.571982438486344 846.7745428400001
	  forward_server_test.go:358: redis mse and avg: 40.478770237088845 794.9892512800001
	  forward_server_test.go:359: mysql mse and avg: 123.02309329266231 1543.54138738
	*/
}

func MSE(a, b float64, x, y []float64) float64 {
	var sum float64
	for i := 0; i < len(x); i++ {
		sum += math.Pow(a+b*x[i]-y[i], 2)
	}
	return math.Pow(sum, 0.5) / float64(len(x))
}
