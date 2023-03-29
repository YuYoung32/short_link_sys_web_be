/**
 * Created by YuYoung on 2023/3/22
 * Description: link handlers
 */

package link

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	. "short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/log"
	"strconv"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var detailsStore []Details

func init() {
	for i := 0; i < MaxSearchAmount; i++ {
		detailsStore = append(detailsStore, Details{
			ShortLink:  RandomString(5) + strconv.Itoa(i),
			LongLink:   "https://baidu.com",
			CreateTime: strconv.FormatInt(time.Now().Unix(), 10),
			Comment:    "无",
		})
	}
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func DetailsListHandler(ctx *gin.Context) {
	ctx.Set("module", "details_handler")
	strSize := ctx.Query("amount")
	strKeyword := ctx.Query("keyword")
	// 必须指定数量
	if strSize == "" {
		ErrMissArgsResp(ctx)
		return
	}
	intSize, err := strconv.Atoi(strSize)
	if err != nil || intSize > MaxSearchAmount {
		ErrInvalidArgsResp(ctx)
		return
	}
	if intSize < 0 {
		intSize = MaxSearchAmount
	}

	var details []Details

	if strKeyword == "" {
		for i := 0; i < MaxSearchAmount && len(details) < intSize; i++ {
			details = append(details, detailsStore[i])
		}
	} else {
		for i := 0; i < MaxSearchAmount && len(details) < intSize; i++ {
			log.MainLogger.WithField("module", "details_handler").Info("keyword: " + strKeyword + ", long_link: " + detailsStore[i].LongLink + ", short_link: " + detailsStore[i].ShortLink)
			if strings.Contains(detailsStore[i].LongLink, strKeyword) || strings.Contains(detailsStore[i].ShortLink, strKeyword) {
				details = append(details, detailsStore[i])
			}
		}
	}

	ctx.JSON(http.StatusOK, details)
}

func AddLinkHandler(ctx *gin.Context) {
	ctx.Set("module", "add_link_handler")
	longLink := ctx.PostForm("longLink")
	if longLink == "" {
		ErrMissArgsResp(ctx)
		return
	}
	log.MainLogger.WithField("module", "add_link_handler").Info("received long_link: " + longLink)
	SuccessGeneralResp(ctx)
}

func DelLinkHandler(ctx *gin.Context) {
	ctx.Set("module", "del_link_handler")
	shortLink := ctx.Query("shortLink")
	if shortLink == "" {
		ErrMissArgsResp(ctx)
		return
	}
	log.MainLogger.WithField("module", "del_link_handler").Info("received short_link: " + shortLink)
	SuccessGeneralResp(ctx)
}

func AmountTotalHandler(ctx *gin.Context) {
	ctx.Set("module", "get_total_amount_handler")
	ctx.JSON(http.StatusOK, AmountTotal{
		Amount: Amount(len(detailsStore)),
	})
}
