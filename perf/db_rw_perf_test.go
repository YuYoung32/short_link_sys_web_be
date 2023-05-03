/**
 * Created by YuYoung on 2023/5/2
 * Description: DB读写测试
 */

package perf

import (
	"github.com/wcharczuk/go-chart"
	"gonum.org/v1/gonum/stat"
	"math/rand"
	"os"
	"short_link_sys_web_be/database"
	"short_link_sys_web_be/link_gen"
	"strconv"
	"testing"
	"time"
)

func TestRW(t *testing.T) {
	//region 数据获取
	amount := 100
	gap := 100
	amountSeq := make([]int, amount)

	random := link_gen.MurmurHash{}
	randomLinks := make([][]database.Link, amount)
	randomWriteRecord := make([]time.Duration, amount)
	randomReadRecord := make([]time.Duration, amount)

	seq := link_gen.SimpleSequencer{}
	seqLinks := make([][]database.Link, amount)
	seqWriteRecord := make([]time.Duration, amount)
	seqReadRecord := make([]time.Duration, amount)

	for i := 0; i < amount; i++ {
		amountSeq[i] = (i + 1) * gap
		for j := 0; j < amountSeq[i]; j++ {
			link := "www.baidu.com" + strconv.Itoa(j)
			randomLinks[i] = append(randomLinks[i], database.Link{
				ShortLink: random.GenLink(link),
			})
			seqLinks[i] = append(seqLinks[i], database.Link{
				ShortLink: seq.GenLink(link),
			})
		}
	}

	db := database.GetDBInstance()
	for i := 0; i < amount; i++ {
		// 随机写
		db.Where("1=1").Delete(&database.Link{})
		start := time.Now()
		db.Create(&randomLinks[i])
		randomWriteRecord[i] = time.Since(start)
		t.Log("random write: ", len(randomLinks[i]), time.Now())
		// 随机读
		// 多次测试, 取平均值
		var sumTime time.Duration
		var cnt = len(randomLinks[i]) / 10
		for j := 0; j < cnt; j++ {
			rand.Seed(time.Now().UnixNano())
			start = time.Now()
			db.Find(&database.Link{
				ShortLink: randomLinks[i][rand.Intn(len(randomLinks[i]))].ShortLink,
			})
			sumTime += time.Since(start)
		}
		randomReadRecord[i] = sumTime / time.Duration(cnt)
		t.Log("random read: ", len(randomLinks[i]), time.Now())
	}

	for i := 0; i < amount; i++ {
		// 顺序写
		db.Where("1=1").Delete(&database.Link{})
		start := time.Now()
		db.Create(&seqLinks[i])
		seqWriteRecord[i] = time.Since(start)
		t.Log("seq write: ", len(seqLinks[i]), time.Now())

		// 顺序读
		var sumTime time.Duration
		var cnt = len(seqLinks[i]) / 10
		for j := 0; j < cnt; j++ {
			rand.Seed(time.Now().UnixNano())
			start = time.Now()
			db.Find(&database.Link{
				ShortLink: seqLinks[i][rand.Intn(len(seqLinks[i]))].ShortLink,
			})
			sumTime += time.Since(start)
		}
		seqReadRecord[i] = sumTime / time.Duration(cnt)
		t.Log("seq read: ", len(seqLinks[i]), time.Now())
	}

	t.Log("随机写入时间：", PrintCommaSeqFloatSlice(DurationsToFloat64(randomWriteRecord)))
	t.Log("随机读取时间：", PrintCommaSeqFloatSlice(DurationsToFloat64(randomReadRecord)))
	t.Log("顺序写入时间：", PrintCommaSeqFloatSlice(DurationsToFloat64(seqWriteRecord)))
	t.Log("顺序读取时间：", PrintCommaSeqFloatSlice(DurationsToFloat64(seqReadRecord)))
	//endregion

	//region 绘图-写入对比
	writeGraph := chart.Chart{
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
			ValueFormatter: func(v interface{}) string {
				i, _ := v.(float64)
				return strconv.Itoa(int(i))
			},
		},

		Series: []chart.Series{
			chart.ContinuousSeries{
				Name: "Random Insertion",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: DurationsToFloat64(randomWriteRecord),
			},
			chart.ContinuousSeries{
				Name: "Sequential Insertion",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(1),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: DurationsToFloat64(seqWriteRecord),
			},
		},
	}
	writeGraph.Elements = []chart.Renderable{
		chart.Legend(&writeGraph),
	}
	f, err := os.Create("perf/db_res/db_write_perf.png")
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
	err = writeGraph.Render(chart.PNG, f)
	if err != nil {
		t.Error(err)
		return
	}
	//endregion

	//region 绘图-读取对比
	readGraph := chart.Chart{
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
			ValueFormatter: func(v interface{}) string {
				i, _ := v.(float64)
				return strconv.Itoa(int(i))
			},
		},

		Series: []chart.Series{
			chart.ContinuousSeries{
				Name: "Random Lookup",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: DurationsToFloat64(randomReadRecord),
			},
			chart.ContinuousSeries{
				Name: "Sequential Lookup",
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(1),
				},
				XValues: IntsToFloat64(amountSeq),
				YValues: DurationsToFloat64(seqReadRecord),
			},
		},
	}
	readGraph.Elements = []chart.Renderable{
		chart.Legend(&readGraph),
	}
	f, err = os.Create("perf/db_res/db_read_perf.png")
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
	err = readGraph.Render(chart.PNG, f)
	if err != nil {
		t.Error(err)
		return
	}
	//endregion
	err = saveGraph(IntsToFloat64(amountSeq),
		DurationsToFloat64(randomWriteRecord),
		DurationsToFloat64(seqWriteRecord),
		"Random Insertion",
		"Sequential Insertion",
		"perf/db_res/db_write_perf.png")
	if err != nil {
		t.Error(err)
		return
	}
	err = saveGraph(IntsToFloat64(amountSeq),
		DurationsToFloat64(randomReadRecord),
		DurationsToFloat64(seqReadRecord),
		"Random Lookup",
		"Sequential Lookup",
		"perf/db_res/db_read_perf.png")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Read Kendall Correlation: ", stat.Kendall(DurationsToFloat64(randomWriteRecord), DurationsToFloat64(seqWriteRecord), nil))
	t.Log("Write Kendall Correlation: ", stat.Kendall(DurationsToFloat64(randomReadRecord), DurationsToFloat64(seqReadRecord), nil))
	/*
	   db_rw_perf_test.go:96: 随机写入时间：[93.7735, 81.6057, 94.1071, 95.377, 93.807, 94.7426, 102.1873, 98.2525, 113.9201, 106.0745, 100.5384, 101.9065, 123.5461, 106.2845, 114.4303, 120.4127, 117.1447, 120.328, 134.4629, 129.6826, 132.8844, 135.6639, 128.9314, 133.6156, 139.6653, 365.645, 382.507, 348.0052, 164.3471, 160.5113, 516.9677, 402.4872, 587.9867, 259.1814, 283.6968, 377.868, 326.2244, 580.8048, 577.8827, 698.8409, 881.3135, 633.8587, 766.8591, 718.4051, 629.221, 709.3574, 848.5789, 943.717, 974.2642, 987.36, 929.2697, 1091.1698, 1110.422, 1078.7387, 1003.0322, 1193.2786, 1071.2724, 1220.0602, 1165.9791, 1205.5705, 1167.6241, 1294.5555, 1298.1654, 1399.6341, 1276.3424, 1439.7327, 1505.344, 1492.3758, 1516.8534, 1601.5786, 1501.6603, 1644.1921, 1747.3891, 1728.7835, 1653.5644, 1705.5479, 1846.3223, 1817.3983, 1855.9381, 1686.74, 1960.8527, 1947.5128, 1796.1359, 1800.015, 1846.6316, 1925.9422, 1953.4948, 2082.9836, 2070.6992, 2091.2007, 2095.4383, 2284.1807, 2255.7378, 2169.8103, 2288.3263, 2289.729, 2472.7911, 2473.3137, 2398.9331, 2573.3472]
	      db_rw_perf_test.go:97: 随机读取时间： [36.59084, 36.22194, 38.7420, 37, 36.628145, 36.430068, 36.45143, 36.903294, 37.82944, 37.487197, 36.686486, 36.496356, 37.398787, 36.577337, 37.763561, 36.977911, 37.139589, 37.550532, 36.479597, 38.354744, 36.766882, 38.21332, 39.82582, 37.170555, 36.995503, 38.157488, 37.19366, 38.433962, 36.920543, 36.842977, 37.558644, 37.02009, 37.015398, 37.232595, 36.5631, 37.049274, 37.06605, 37.374692, 37.09022, 37.473225, 40.191575, 37.048906, 38.395618, 38.05657, 37.234285, 36.954417, 37.121915, 37.439588, 37.486124, 37.283635, 37.340365, 37.385613, 37.73604, 37.031308, 37.152394, 37.809976, 36.929506, 36.89862, 36.930042, 37.033472, 36.801815, 36.908769, 36.945722, 36.905932, 37.146937, 36.96288, 36.848517, 37.058402, 40.463663, 37.350487, 38.580231, 36.7521, 36.718097, 37.346725, 36.759009, 36.822814, 37.091716, 37.176208, 37.075517, 37.360643, 38.632622, 37.011674, 37.08949, 37.062728, 36.942513, 37.202516, 36.740618, 37.120297, 38.561081, 37.204442, 37.295326, 37.034373, 36.962561, 36.996266, 36.797883, 37.258847, 36.926255, 36.85416, 36.722086, 36.846818, 37.020016]
	      db_rw_perf_test.go:98: 顺序写入时间： [102.1173, 119.888, 137.0565, 147.6728, 178.7318, 175.6819, 178.5351, 183.8397, 167.995, 171.8231, 186.2256, 177.1653, 182.2785, 189.0089, 179.531, 204.09, 190.133, 195.221, 192.3728, 204.2451, 202.61, 195.6841, 200.6804, 203.4707, 198.5648, 323.0315, 315.8626, 366.656, 342.8883, 394.0232, 436.5154, 512.0596, 517.0621, 592.5516, 465.9549, 656.5976, 595.2272, 415.9135, 355.1909, 574.1063, 692.7143, 792.8948, 705.2353, 795.6493, 829.7511, 842.6349, 954.5081, 989.0082, 925.4198, 974.544, 1066.0374, 1084.9307, 1150.1058, 981.2822, 1161.0601, 1075.454, 1053.0071, 1179.0119, 1114.0664, 1141.96, 1238.1853, 1219.9752, 1367.1043, 1429.1906, 1317.9494, 1431.6896, 1442.6032, 1412.9878, 1540.6499, 1400.6862, 1443.1047, 1579.5414, 1559.5123, 1662.1197, 1681.7168, 1782.7878, 1767.0886, 1708.6926, 1873.7523, 1862.3001, 1895.2474, 1935.7247, 1938.0334, 2079.9704, 1814.0897, 1919.2393, 1975.6826, 1977.8594, 2145.3982, 2031.3946, 2223.3684, 2314.7398, 2162.5544, 2236.8212, 2255.2697, 2403.2279, 2278.4615, 2347.1759, 2270.315, 2504.1725]
	      db_rw_perf_test.go:99: 顺序读取时间： [39.76726, 36.926715, 37.103, 36, 37.753557, 37.278934, 37.130736, 37.362395, 37.439585, 37.443375, 38.081296, 37.353861, 37.590791, 37.392025, 36.939486, 37.149344, 37.248281, 37.11503, 36.72955, 37.722144, 36.631337, 36.738257, 36.753856, 37.090015, 36.884286, 42.840307, 36.993568, 38.078036, 36.871817, 36.819918, 36.894138, 37.198352, 37.102216, 36.990859, 37.066573, 37.075994, 37.124887, 37.624682, 36.848047, 36.815727, 37.075607, 37.083764, 37.542695, 37.157149, 37.481458, 37.027168, 37.233472, 37.052891, 37.332093, 37.297038, 37.205943, 36.921381, 37.20092, 37.084049, 36.961999, 37.004671, 36.997374, 36.916034, 36.949042, 37.506475, 37.02483, 39.53264, 37.482549, 37.109747, 37.797462, 36.933733, 37.043365, 37.010789, 37.727726, 36.769263, 36.580314, 37.260127, 36.94711, 37.004624, 37.590562, 36.814087, 37.265906, 36.844409, 36.77119, 37.617189, 36.96516, 37.639397, 38.904803, 37.462577, 37.634942, 37.263588, 37.137993, 36.886297, 37.041715, 37.133769, 45.017187, 39.758264, 38.694891, 38.867519, 39.63058, 40.56498, 39.355863, 39.277726, 41.407952, 38.858639, 39.554037]
	      db_rw_perf_test.go:228: Read Kendall Correlation:  0.92
	      db_rw_perf_test.go:229: Write Kendall Correlation:  -0.053737373737373736
	*/

}

func saveGraph(x, y1, y2 []float64, y1Name, y2Name, path string) error {
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

func TestProcessData(t *testing.T) {
	amountSeq := make([]int, 100)
	for i := 0; i < 100; i++ {
		amountSeq[i] = (i + 1) * 100
	}
	randomWriteRecord := []float64{93.7735, 81.6057, 94.1071, 95.377, 93.807, 94.7426, 102.1873, 98.2525, 113.9201, 106.0745, 100.5384, 101.9065, 123.5461, 106.2845, 114.4303, 120.4127, 117.1447, 120.328, 134.4629, 129.6826, 132.8844, 135.6639, 128.9314, 133.6156, 139.6653, 365.645, 382.507, 348.0052, 164.3471, 160.5113, 516.9677, 402.4872, 587.9867, 259.1814, 283.6968, 377.868, 326.2244, 580.8048, 577.8827, 698.8409, 881.3135, 633.8587, 766.8591, 718.4051, 629.221, 709.3574, 848.5789, 943.717, 974.2642, 987.36, 929.2697, 1091.1698, 1110.422, 1078.7387, 1003.0322, 1193.2786, 1071.2724, 1220.0602, 1165.9791, 1205.5705, 1167.6241, 1294.5555, 1298.1654, 1399.6341, 1276.3424, 1439.7327, 1505.344, 1492.3758, 1516.8534, 1601.5786, 1501.6603, 1644.1921, 1747.3891, 1728.7835, 1653.5644, 1705.5479, 1846.3223, 1817.3983, 1855.9381, 1686.74, 1960.8527, 1947.5128, 1796.1359, 1800.015, 1846.6316, 1925.9422, 1953.4948, 2082.9836, 2070.6992, 2091.2007, 2095.4383, 2284.1807, 2255.7378, 2169.8103, 2288.3263, 2289.729, 2472.7911, 2473.3137, 2398.9331, 2573.3472}
	randomReadRecord := []float64{36.59084, 36.22194, 38.7420, 37, 36.628145, 36.430068, 36.45143, 36.903294, 37.82944, 37.487197, 36.686486, 36.496356, 37.398787, 36.577337, 37.763561, 36.977911, 37.139589, 37.550532, 36.479597, 38.354744, 36.766882, 38.21332, 39.82582, 37.170555, 36.995503, 38.157488, 37.19366, 38.433962, 36.920543, 36.842977, 37.558644, 37.02009, 37.015398, 37.232595, 36.5631, 37.049274, 37.06605, 37.374692, 37.09022, 37.473225, 40.191575, 37.048906, 38.395618, 38.05657, 37.234285, 36.954417, 37.121915, 37.439588, 37.486124, 37.283635, 37.340365, 37.385613, 37.73604, 37.031308, 37.152394, 37.809976, 36.929506, 36.89862, 36.930042, 37.033472, 36.801815, 36.908769, 36.945722, 36.905932, 37.146937, 36.96288, 36.848517, 37.058402, 40.463663, 37.350487, 38.580231, 36.7521, 36.718097, 37.346725, 36.759009, 36.822814, 37.091716, 37.176208, 37.075517, 37.360643, 38.632622, 37.011674, 37.08949, 37.062728, 36.942513, 37.202516, 36.740618, 37.120297, 38.561081, 37.204442, 37.295326, 37.034373, 36.962561, 36.996266, 36.797883, 37.258847, 36.926255, 36.85416, 36.722086, 36.846818, 37.020016}
	seqWriteRecord := []float64{102.1173, 119.888, 137.0565, 147.6728, 178.7318, 175.6819, 178.5351, 183.8397, 167.995, 171.8231, 186.2256, 177.1653, 182.2785, 189.0089, 179.531, 204.09, 190.133, 195.221, 192.3728, 204.2451, 202.61, 195.6841, 200.6804, 203.4707, 198.5648, 323.0315, 315.8626, 366.656, 342.8883, 394.0232, 436.5154, 512.0596, 517.0621, 592.5516, 465.9549, 656.5976, 595.2272, 415.9135, 355.1909, 574.1063, 692.7143, 792.8948, 705.2353, 795.6493, 829.7511, 842.6349, 954.5081, 989.0082, 925.4198, 974.544, 1066.0374, 1084.9307, 1150.1058, 981.2822, 1161.0601, 1075.454, 1053.0071, 1179.0119, 1114.0664, 1141.96, 1238.1853, 1219.9752, 1367.1043, 1429.1906, 1317.9494, 1431.6896, 1442.6032, 1412.9878, 1540.6499, 1400.6862, 1443.1047, 1579.5414, 1559.5123, 1662.1197, 1681.7168, 1782.7878, 1767.0886, 1708.6926, 1873.7523, 1862.3001, 1895.2474, 1935.7247, 1938.0334, 2079.9704, 1814.0897, 1919.2393, 1975.6826, 1977.8594, 2145.3982, 2031.3946, 2223.3684, 2314.7398, 2162.5544, 2236.8212, 2255.2697, 2403.2279, 2278.4615, 2347.1759, 2270.315, 2504.1725}
	seqReadRecord := []float64{39.76726, 36.926715, 37.103, 36, 37.753557, 37.278934, 37.130736, 37.362395, 37.439585, 37.443375, 38.081296, 37.353861, 37.590791, 37.392025, 36.939486, 37.149344, 37.248281, 37.11503, 36.72955, 37.722144, 36.631337, 36.738257, 36.753856, 37.090015, 36.884286, 42.840307, 36.993568, 38.078036, 36.871817, 36.819918, 36.894138, 37.198352, 37.102216, 36.990859, 37.066573, 37.075994, 37.124887, 37.624682, 36.848047, 36.815727, 37.075607, 37.083764, 37.542695, 37.157149, 37.481458, 37.027168, 37.233472, 37.052891, 37.332093, 37.297038, 37.205943, 36.921381, 37.20092, 37.084049, 36.961999, 37.004671, 36.997374, 36.916034, 36.949042, 37.506475, 37.02483, 39.53264, 37.482549, 37.109747, 37.797462, 36.933733, 37.043365, 37.010789, 37.727726, 36.769263, 36.580314, 37.260127, 36.94711, 37.004624, 37.590562, 36.814087, 37.265906, 36.844409, 36.77119, 37.617189, 36.96516, 37.639397, 38.904803, 37.462577, 37.634942, 37.263588, 37.137993, 36.886297, 37.041715, 37.133769, 45.017187, 39.758264, 38.694891, 38.867519, 39.63058, 40.56498, 39.355863, 39.277726, 41.407952, 38.858639, 39.554037}
	err := saveGraph(IntsToFloat64(amountSeq),
		randomWriteRecord,
		seqWriteRecord,
		"Random Insertion",
		"Sequential Insertion",
		"perf/db_res/db_write_perf.png")
	if err != nil {
		t.Error(err)
	}
	err = saveGraph(IntsToFloat64(amountSeq),
		randomReadRecord,
		seqReadRecord,
		"Random Lookup",
		"Sequential Lookup",
		"perf/db_res/db_read_perf.png")
	if err != nil {
		t.Error(err)
	}

	t.Log("Read Kendall Correlation: ", stat.Kendall(randomWriteRecord, seqWriteRecord, nil))
	t.Log("Write Kendall Correlation: ", stat.Kendall(randomReadRecord, seqReadRecord, nil))
}

func TestCorrelation(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5}
	b := []float64{1, 2, 3, 4, 5}
	t.Log(stat.Kendall(a, b, nil))
}
