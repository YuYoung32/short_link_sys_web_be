/**
 * Created by YuYoung on 2023/4/21
 * Description: 登录
 */

package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/log"
	"short_link_sys_web_be/utils"
	"time"
)

func LoginHandler(ctx *gin.Context) {
	logger := log.GetLogger()
	password := conf.GlobalConfig.GetString("auth.password")
	unAuthedPassword := ctx.Query("password")

	if password != unAuthedPassword {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "口令错误",
		})
		return
	}
	token, err := utils.GenerateToken()
	if err != nil {
		logger.Error("token生成失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "token生成失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token":  token,
		"expire": time.Now().Add(utils.Expire).UnixMilli(),
	})
}
