/**
 * Created by YuYoung on 2023/3/22
 * Description: link handlers
 */

package link

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/database"
	. "short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/link_gen"
	"short_link_sys_web_be/log"
	"strconv"
	"strings"
)

var linkGenAlgorithm link_gen.LinkGen

func Init() {
	mapAlgorithm := map[string]link_gen.LinkGen{
		"murmurHash":   link_gen.MurmurHash{},
		"xxHash":       link_gen.XXHash{},
		"fnvHash":      link_gen.FNVHash{},
		"simpleSeq":    link_gen.SimpleSequencer{},
		"snowflakeSeq": link_gen.SnowflakeSequencer{},
	}

	if linkGenAlgorithm = mapAlgorithm[conf.GlobalConfig.GetString("handler.link.algorithm")]; linkGenAlgorithm == nil {
		linkGenAlgorithm = mapAlgorithm["simpleSeq"]
		log.GetLogger().Error("bad config: handler.link.algorithm")
		return
	}
	log.GetLogger().Info("set link algorithm to: " + conf.GlobalConfig.GetString("handler.link.algorithm"))
}

func DetailsListHandler(ctx *gin.Context) {
	db := database.GetDBInstance()

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
	db.Where(strings.Join(queryTemplateList, " or "), queryArgsList...).Limit(intSize).Find(&details.Links)
	ctx.JSON(http.StatusOK, details)
}

func AddLinkHandler(ctx *gin.Context) {
	db := database.GetDBInstance()

	var queryAddListBind []struct {
		LongLink string `json:"longLink"`
		Comment  string `json:"comment"`
	}
	if err := ctx.BindJSON(&queryAddListBind); err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}

	var detailsStore []database.Link
	for i, link := range queryAddListBind {
		detailsStore = append(detailsStore, database.Link{
			ShortLink: "sgews" + strconv.Itoa(i),
			LongLink:  link.LongLink,
			Comment:   link.Comment,
		})
	}

	db.Create(&detailsStore)
	SuccessGeneralResp(ctx)
}

func DelLinkHandler(ctx *gin.Context) {
	ctx.Set("module", "del_link_handler")
	db := database.GetDBInstance()

	var queryDelListBind struct {
		ShortLinks []string `json:"shortLinks"`
	}
	if err := ctx.BindJSON(&queryDelListBind); err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}
	for _, shortLink := range queryDelListBind.ShortLinks {
		db.Delete(&database.Link{}, "short_link = ?", shortLink)
	}
	SuccessGeneralResp(ctx)
}

func UpdateLinkHandler(ctx *gin.Context) {
	db := database.GetDBInstance()

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
	db.First(&link)
	link.ShortLink = "newsl"
	link.LongLink = queryUpdateListBind.LongLink
	link.Comment = queryUpdateListBind.Comment
	db.Create(&link)
	SuccessGeneralResp(ctx)
}

func AmountTotalHandler(ctx *gin.Context) {
	db := database.GetDBInstance()

	var amount AmountTotal
	db.Model(&database.Link{}).Count(&amount.Amount)
	ctx.JSON(http.StatusOK, amount)
}
