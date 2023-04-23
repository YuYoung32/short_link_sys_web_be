/**
 * Created by YuYoung on 2023/4/21
 * Description: 权限认证
 */

package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/utils"
)

func AuthMiddleware(ctx *gin.Context) {
	var token string
	reqAuth := ctx.Request.Header.Get("Authorization")
	if reqAuth == "" {
		token = ctx.Query("token")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "no Authorization field",
			})
			ctx.Abort()
			return
		}
	} else {
		s := bytes.Split([]byte(reqAuth), []byte(" "))
		if len(s) != 2 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "no Authorization field",
			})
			ctx.Abort()
			return
		}
		token = string(s[1])
	}

	ok, err := utils.ValidToken(token)
	if err != nil {
		log.GetLogger().Debug("token error: ", err)
		ctx.Abort()
		return
	}
	if !ok {
		log.GetLogger().Debug("token invalid: ", token)
		ctx.Abort()
		return
	}

	ctx.Next()
}
