/**
 * Created by YuYoung on 2023/3/22
 * Description: link handlers
 */

package link

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"net/url"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/database"
	. "short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/link_gen"
	"short_link_sys_web_be/log"
	"strconv"
	"strings"
	"time"
)

var linkGenAlgorithm link_gen.LinkGen
var shortLinkBF *bloom.BloomFilter
var longLinkBF *bloom.BloomFilter

func init() {
	logger := log.GetLogger()
	mapAlgorithm := map[string]link_gen.LinkGen{
		"murmurHash":   link_gen.MurmurHash{},
		"xxHash":       link_gen.XXHash{},
		"fnvHash":      link_gen.FNVHash{},
		"simpleSeq":    link_gen.SimpleSequencer{},
		"snowflakeSeq": link_gen.SnowflakeSequencer{},
	}

	algorithmName := conf.GlobalConfig.GetString("handler.link.algorithm")
	if linkGenAlgorithm = mapAlgorithm[algorithmName]; linkGenAlgorithm == nil {
		linkGenAlgorithm = mapAlgorithm["simpleSeq"]
		logger.Error("bad config: handler.link.algorithm")
		return
	} else {
		logger.Info("set link algorithm to: " + conf.GlobalConfig.GetString("handler.link.algorithm"))
	}

	longLinkBF = bloomFilterInit("long_link")
	// 使用hash算法时, 才会使用短链的布隆过滤器
	if linkGenAlgorithm.GetType() == link_gen.HashType {
		shortLinkBF = bloomFilterInit("short_link")
	}
	go func() {
		for {
			time.Sleep(time.Hour * 24)
			longLinkBF = bloomFilterInit("long_link")
			if linkGenAlgorithm.GetType() == link_gen.HashType {
				shortLinkBF = bloomFilterInit("short_link")
			}
		}
	}()
}

// bloomFilterInit 布隆过滤器创造与初始化
func bloomFilterInit(key string) *bloom.BloomFilter {
	preFix := "handler.link.bloomFilter."
	falsePositiveRateName := preFix + "falsePositiveRate"
	expectedNumberOfElementsName := preFix + "expectedNumberOfElements"
	needToLoadName := preFix + "needToLoad"
	if !(conf.GlobalConfig.IsSet(falsePositiveRateName) &&
		conf.GlobalConfig.IsSet(expectedNumberOfElementsName) &&
		conf.GlobalConfig.IsSet(needToLoadName)) {
		log.GetLoggerWithSkip(2).Error("bad config: " + falsePositiveRateName + " or " + expectedNumberOfElementsName)
	}
	bf := bloom.NewWithEstimates(
		conf.GlobalConfig.GetUint(expectedNumberOfElementsName),
		conf.GlobalConfig.GetFloat64(falsePositiveRateName))

	if conf.GlobalConfig.GetBool(needToLoadName) {
		now := time.Now()
		log.GetLoggerWithSkip(2).Info("load data to bloom filter... ", now.Format("2006-01-02 15:04:05"))
		db := database.GetMysqlInstance()
		var links []string
		db.Model(&database.Link{}).Pluck(key, &links)
		for _, link := range links {
			bf.AddString(link)
		}
		after := time.Now()
		log.GetLoggerWithSkip(2).Info("succeed to load data to bloom filter, uses ", after.Sub(now).Seconds(), "s")
	}
	return bf
}

func genUniqueLink(longLink string) string {
	var shortLink string
	switch linkGenAlgorithm.GetType() {
	case link_gen.HashType:
		for {
			shortLink = linkGenAlgorithm.GenLink(longLink)
			if !shortLinkBF.TestString(shortLink) {
				break
			}
			longLink = longLink + strconv.Itoa(rand.Int())
			shortLink = linkGenAlgorithm.GenLink(longLink)
		}
	case link_gen.SeqType:
		shortLink = linkGenAlgorithm.GenLink(longLink)
	}
	return shortLink
}

func DetailsListHandler(ctx *gin.Context) {
	db := database.GetMysqlInstance()

	var intSize int
	var queryTemplateList []string
	var queryArgsList []interface{}

	if strSize := ctx.Query("amount"); strSize == "" {
		ErrMissArgsResp(ctx)
		return
	} else {
		var err error
		if intSize, err = strconv.Atoi(strSize); err != nil || intSize > MaxSearchAmount {
			ErrInvalidArgsResp(ctx)
			return
		}
		if intSize < 0 {
			intSize = MaxSearchAmount
		}
	}
	if shortLink := ctx.Query("shortLink"); shortLink != "" {
		queryTemplateList = append(queryTemplateList, "short_link like (?)")
		queryArgsList = append(queryArgsList, "%"+shortLink+"%")
	}
	if longLink := ctx.Query("longLink"); longLink != "" {
		queryTemplateList = append(queryTemplateList, "long_link like (?)")
		queryArgsList = append(queryArgsList, "%"+longLink+"%")
	}
	if comment := ctx.Query("comment"); comment != "" {
		queryTemplateList = append(queryTemplateList, "comment like (?)")
		queryArgsList = append(queryArgsList, "%"+comment+"%")
	}

	var details DetailsListResponse
	db.Model(&database.Link{}).Where(strings.Join(queryTemplateList, " or "), queryArgsList...).Count(&details.LinksTotal)
	db.Where(strings.Join(queryTemplateList, " or "), queryArgsList...).Limit(intSize).Order("update_time desc").Find(&details.Links)
	ctx.JSON(http.StatusOK, details)
}

func processRawLink(rawLink string) string {
	parse, err := url.Parse(rawLink)
	if err != nil {
		log.GetLogger().Error("url parse error: ", err)
	}
	if parse.Scheme != "" && parse.Host != "" {
		return rawLink
	}
	return "http://" + rawLink
}

func AddLinkHandler(ctx *gin.Context) {
	db := database.GetMysqlInstance()

	var queryAddListBind []struct {
		LongLink string `json:"longLink"`
		Comment  string `json:"comment"`
	}
	if err := ctx.BindJSON(&queryAddListBind); err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}

	var detailsStore []database.Link
	for _, link := range queryAddListBind {
		// 绝对连接检测与转换
		longLink := processRawLink(link.LongLink)
		// 布隆过滤器检测是否已经存在长链, 若存在无需再次添加
		if longLinkBF.TestString(longLink) {
			// 排除假阳性
			err := db.Where(&database.Link{LongLink: longLink}).Take(&database.Link{}).Error
			if err == nil {
				// 已确认数据库中有该长链
				continue
			} else if err != gorm.ErrRecordNotFound {
				log.GetLogger().Error(err)
				return
			}
		}

		// 已确认数据库中无长链
		longLinkBF.AddString(longLink)
		shortLink := genUniqueLink(longLink)
		if shortLinkBF != nil {
			shortLinkBF.AddString(shortLink)
		}
		detailsStore = append(detailsStore, database.Link{
			ShortLink: shortLink,
			LongLink:  longLink,
			Comment:   link.Comment,
		})
	}

	db.Create(detailsStore)
	SuccessGeneralResp(ctx)
}

func DelLinkHandler(ctx *gin.Context) {
	ctx.Set("module", "del_link_handler")
	db := database.GetMysqlInstance()

	var queryDelListBind struct {
		ShortLinks []string `json:"shortLinks"`
	}
	if err := ctx.BindJSON(&queryDelListBind); err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}
	db.Where("short_link in (?)", queryDelListBind.ShortLinks).Delete(&database.Link{})
	SuccessGeneralResp(ctx)
}

func UpdateLinkHandler(ctx *gin.Context) {
	db := database.GetMysqlInstance()

	var queryUpdateListBind struct {
		ShortLink string `json:"shortLink"`
		LongLink  string `json:"longLink"`
		Comment   string `json:"comment"`
	}
	if err := ctx.BindJSON(&queryUpdateListBind); err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}

	var link database.Link
	link.ShortLink = queryUpdateListBind.ShortLink
	db.Take(&link)
	link.LongLink = queryUpdateListBind.LongLink
	link.Comment = queryUpdateListBind.Comment
	db.Updates(&link)
	SuccessGeneralResp(ctx)
}

func AmountTotalHandler(ctx *gin.Context) {
	db := database.GetMysqlInstance()

	var amount AmountTotal
	db.Model(&database.Link{}).Count(&amount.Amount)
	ctx.JSON(http.StatusOK, amount)
}
