/**
 * Created by YuYoung on 2023/3/22
 * Description: link handlers
 */

package link

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "short_link_sys_web_be/handler/common"
	"short_link_sys_web_be/log"
	"strconv"
)

func testDetailsDataGenerator(size int) []Details {
	var details []Details
	for i := 0; i < size; i++ {
		details = append(details, Details{
			ShortLink:  "http://localhost:8080/shortLink",
			LongLink:   "http://localhost:8080/longLink",
			CreateTime: "1679533954",
		})
	}
	return details
}

func DetailsListHandler(ctx *gin.Context) {
	ctx.Set("module", "details_handler")
	strSize := ctx.Query("amount")
	if strSize == "" {
		ErrMissArgsResp(ctx)
		return
	}
	intSize, err := strconv.Atoi(strSize)
	if err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}
	if intSize < 0 || intSize > MaxSearchAmount {
		ErrInvalidArgsResp(ctx)
		return
	}

	details := testDetailsDataGenerator(intSize)
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
		Amount: Amount(100),
	})
}
