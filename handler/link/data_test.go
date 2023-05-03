/**
 * Created by YuYoung on 2023/4/25
 * Description: 生成测试访问记录
 */

package link

import (
	"math/rand"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/database"
	"short_link_sys_web_be/link_gen"
	"short_link_sys_web_be/log"
	"strconv"
	"testing"
	"time"
)

func TestGenerateLinkData(t *testing.T) {
	conf.GlobalConfig.SetConfigName("config")
	conf.GlobalConfig.SetConfigType("yaml")
	conf.GlobalConfig.AddConfigPath("../../conf")
	if err := conf.GlobalConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	dataList := []database.Link{
		{LongLink: "http://www.taobao.com", Comment: "测试1"},
		{LongLink: "http://jd.com", Comment: "测试3"},
		{LongLink: "http://www.computer.hdu.edu.cn", Comment: "测试2"},
		{LongLink: "https://www.baidu.com/", Comment: "百度一下你就知道"},
		{LongLink: "https://www.google.com/", Comment: "谷歌搜索引擎"},
		{LongLink: "https://www.microsoft.com/zh-cn", Comment: "微软官网"},
		{LongLink: "https://www.apple.com/cn/", Comment: "苹果中国官网"},
		{LongLink: "https://www.qq.com/", Comment: "腾讯首页"},
		{LongLink: "https://www.sina.com.cn/", Comment: "新浪首页"},
		{LongLink: "https://www.163.com/", Comment: "网易首页"},
		{LongLink: "https://www.sohu.com/", Comment: "搜狐首页"},
		{LongLink: "https://www.vip.com/", Comment: "唯品会官网"},
		{LongLink: "https://www.suning.com/", Comment: "苏宁易购官网"},
		{LongLink: "https://www.jd.com/", Comment: "京东官网"},
		{LongLink: "https://www.douban.com/", Comment: "豆瓣"},
		{LongLink: "https://www.zhihu.com/", Comment: "知乎"},
		{LongLink: "https://www.cnblogs.com/", Comment: "博客园"},
		{LongLink: "http://www.people.com.cn/", Comment: "人民网"},
		{LongLink: "https://www.xinhuanet.com/", Comment: "新华网"},
		{LongLink: "http://www.chinadaily.com.cn/", Comment: "中国日报"},
		{LongLink: "https://www.ftchinese.com/", Comment: "FT中文网"},
		{LongLink: "https://www.thepaper.cn/", Comment: "澎湃新闻"},
		{LongLink: "https://www.cctv.com/", Comment: "央视网"},
		{LongLink: "http://www.china.com.cn/", Comment: "中国网"},
		{LongLink: "https://www.toutiao.com/", Comment: "今日头条"},
		{LongLink: "https://www.bilibili.com/", Comment: "哔哩哔哩"},
		{LongLink: "https://www.acfun.cn/", Comment: "AcFun弹幕视频网"},
		{LongLink: "https://www.youtube.com/", Comment: "YouTube"},
		{LongLink: "https://www.netflix.com/cn/", Comment: "Netflix"},
		{LongLink: "https://www.amazon.cn/", Comment: "亚马逊中国"},
		{LongLink: "https://www.tmall.com/", Comment: "天猫"},
		{LongLink: "https://www.alibaba.com/", Comment: "阿里巴巴国际站"},
		{LongLink: "https://www.huawei.com/cn/", Comment: "华为官网"},
		{LongLink: "https://www.mi.com/", Comment: "小米官网"},
		{LongLink: "https://www.oppo.com/cn/", Comment: "OPPO官网"},
		{LongLink: "https://www.vivo.com/cn/", Comment: "vivo官网"},
		{LongLink: "https://www.realme.com/cn/", Comment: "realme官网"},
		{LongLink: "https://www.oneplus.com/cn/", Comment: "一加手机官网"},
		{LongLink: "https://www.letv.com/", Comment: "乐视视频"},
		{LongLink: "https://v.qq.com/", Comment: "腾讯视频"},
		{LongLink: "https://www.youku.com/", Comment: "优酷视频"},
		{LongLink: "https://www.iqiyi.com/", Comment: "爱奇艺视频"},
		{LongLink: "http://www.tmtpost.com/", Comment: "钛媒体"},
		{LongLink: "https://36kr.com/", Comment: "36氪"},
		{LongLink: "https://www.ifanr.com/", Comment: "爱范儿"},
		{LongLink: "https://sspai.com/", Comment: "少数派"},
		{LongLink: "https://www.zdnet.com/cn/", Comment: "ZDNet中国"},
		{LongLink: "http://www.dzwww.com/", Comment: "大众网"},
		{LongLink: "https://www.guancha.cn/", Comment: "观察者网"},
		{LongLink: "https://www.jiemian.com/", Comment: "界面新闻"},
		{LongLink: "https://www.huxiu.com/", Comment: "虎嗅网"},
		{LongLink: "https://www.woshipm.com/", Comment: "人人都是产品经理"},
		{LongLink: "https://www.oeeee.com/", Comment: "南方周末"},
		{LongLink: "https://www.thepaper.cn/", Comment: "澎湃新闻"},
		{LongLink: "https://www.csdn.net/", Comment: "CSDN技术社区"},
		{LongLink: "https://www.oschina.net/", Comment: "开源中国"},
		{LongLink: "http://www.chinanews.com/", Comment: "中国新闻网"},
		{LongLink: "https://www.huanqiu.com/", Comment: "环球网"},
		{LongLink: "http://www.xinhuanet.com/politics/", Comment: "新华网-时政频道"},
		{LongLink: "https://www.theguardian.com/international", Comment: "The Guardian"},
		{LongLink: "https://www.bbc.com/news/world", Comment: "BBC News"},
		{LongLink: "https://edition.cnn.com/", Comment: "CNN"},
		{LongLink: "https://www.reuters.com/world", Comment: "Reuters"},
		{LongLink: "https://www.nytimes.com/", Comment: "The New York Times"},
		{LongLink: "https://www.washingtonpost.com/", Comment: "The Washington Post"},
		{LongLink: "https://www.lemonde.fr/", Comment: "Le Monde"},
		{LongLink: "https://www.spiegel.de/", Comment: "Der Spiegel"},
		{LongLink: "https://www.economist.com/", Comment: "The Economist"},
		{LongLink: "https://www.ft.com/", Comment: "Financial Times"},
		{LongLink: "https://www.wsj.com/", Comment: "The Wall Street Journal"},
		{LongLink: "https://www.thetimes.co.uk/", Comment: "The Times"},
		{LongLink: "https://www.scmp.com/", Comment: "South China Morning Post"},
	}

	log.GetLogger().Info("will generate", len(dataList), "item")

	err := GenerateLinkData(dataList)
	if err != nil {
		t.Error(err)
		return
	}

	link_gen.Terminate()
}

func TestGenerateVisitData(t *testing.T) {
	// 配置文件初始化
	conf.GlobalConfig.SetConfigName("config")
	conf.GlobalConfig.SetConfigType("yaml")
	conf.GlobalConfig.AddConfigPath("../../conf")
	if err := conf.GlobalConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	db := database.GetDBInstance()

	var shortLinks []string
	type Visit struct {
		ShortLink string `json:"shortLink" gorm:"type:varchar(255) COLLATE utf8_bin"`
		IP        string `json:"ip"`
		Region    string `json:"region"`
		VisitTime int64  `json:"visitTime" gorm:"autoCreateTime"`
	}
	db.Model(&database.Link{}).Pluck("short_link", &shortLinks)
	shortLinksLen := len(shortLinks)
	provinces := []string{
		"北京市",
		"天津市",
		"河北省",
		"山西省",
		"内蒙古自治区",
		"辽宁省",
		"吉林省",
		"黑龙江省",
		"上海市",
		"江苏省",
		"浙江省",
		"安徽省",
		"福建省",
		"江西省",
		"山东省",
		"河南省",
		"湖北省",
		"湖南省",
		"广东省",
		"广西壮族自治区",
		"海南省",
		"重庆市",
		"四川省",
		"贵州省",
		"云南省",
		"西藏自治区",
		"陕西省",
		"甘肃省",
		"青海省",
		"宁夏回族自治区",
		"新疆维吾尔自治区",
		"台湾省",
		"香港特别行政区",
		"澳门特别行政区",
	}
	var beginTS int64 = 1681488000 //2023-04-15 00:00:00 +0800 CST
	var cross = time.Now().Unix() - beginTS
	var addList []Visit
	for i := 0; i < 1000; i++ {
		rand.Seed(int64(i))
		addList = append(addList, Visit{
			ShortLink: shortLinks[rand.Intn(shortLinksLen)],
			IP:        strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)),
			Region:    provinces[rand.Intn(len(provinces))],
			VisitTime: rand.Int63n(cross) + beginTS,
		})
	}
	db.Model(&Visit{}).Create(&addList)
}
