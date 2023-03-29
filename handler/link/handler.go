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
		for i := 0; i < len(detailsStore) && len(details) < intSize; i++ {
			details = append(details, detailsStore[i])
		}
	} else {
		for i := 0; i < len(detailsStore) && len(details) < intSize; i++ {
			if strings.Contains(detailsStore[i].LongLink, strKeyword) || strings.Contains(detailsStore[i].ShortLink, strKeyword) {
				details = append(details, detailsStore[i])
			}
		}
	}

	ctx.JSON(http.StatusOK, details)
}

func AddLinkHandler(ctx *gin.Context) {
	ctx.Set("module", "add_link_handler")
	var addLinkList []struct {
		LongLink string `json:"longLink"`
		Comment  string `json:"comment"`
	}
	time.Sleep(1 * time.Second)
	if err := ctx.BindJSON(&addLinkList); err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}
	for i, link := range addLinkList {
		detailsStore = append(detailsStore, Details{
			ShortLink:  RandomString(5) + strconv.Itoa(i),
			LongLink:   link.LongLink,
			CreateTime: strconv.FormatInt(time.Now().Unix(), 10),
			Comment:    link.Comment,
		})

	}
	SuccessGeneralResp(ctx)
}

func DelLinkHandler(ctx *gin.Context) {
	ctx.Set("module", "del_link_handler")
	var shortLinks struct {
		ShortLinks []string `json:"shortLinks"`
	}
	if err := ctx.BindJSON(&shortLinks); err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}
	for _, shortLink := range shortLinks.ShortLinks {
		for i, detail := range detailsStore {
			if detail.ShortLink == shortLink {
				detailsStore = append(detailsStore[:i], detailsStore[i+1:]...)
			}
		}
	}
	SuccessGeneralResp(ctx)
}

func UpdateLinkHandler(ctx *gin.Context) {
	ctx.Set("module", "update_link_handler")
	var updateLink struct {
		ShortLink string `json:"shortLink"`
		LongLink  string `json:"longLink"`
		Comment   string `json:"comment"`
	}
	if err := ctx.BindJSON(&updateLink); err != nil {
		ErrInvalidArgsResp(ctx)
		return
	}
	for i, detail := range detailsStore {
		if detail.ShortLink == updateLink.ShortLink {
			detailsStore[i].LongLink = updateLink.LongLink
			detailsStore[i].Comment = updateLink.Comment
		}
	}
	SuccessGeneralResp(ctx)
}

func AmountTotalHandler(ctx *gin.Context) {
	ctx.Set("module", "get_total_amount_handler")
	ctx.JSON(http.StatusOK, AmountTotal{
		Amount: len(detailsStore),
	})
}
