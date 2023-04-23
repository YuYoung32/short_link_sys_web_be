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
	reqAuth := ctx.Request.Header.Get("Authorization")
	s := bytes.Split([]byte(reqAuth), []byte(" "))
	if len(s) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "请求头中无Authorization字段或错误格式",
		})
		ctx.Abort()
		return
	}

	ok, err := utils.ValidToken(string(s[1]))
	if err != nil {
		log.GetLogger().Debug("token error: ", err)
		ctx.Abort()
		return
	}
	if !ok {
		log.GetLogger().Debug("token invalid: ", string(s[1]))
		ctx.Abort()
		return
	}

	ctx.Next()
}
