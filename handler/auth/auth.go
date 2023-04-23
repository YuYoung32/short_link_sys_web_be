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
	var bind struct {
		Password string `json:"password" binding:"required"`
	}
	err := ctx.ShouldBind(&bind)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if password != bind.Password {
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

func ChangePasswordHandler(ctx *gin.Context) {
	var bind struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	err := ctx.ShouldBind(&bind)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if bind.OldPassword != conf.GlobalConfig.GetString("auth.password") {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "wrong password",
		})
		return
	}
	conf.GlobalConfig.Set("auth.password", bind.NewPassword)
	err = conf.GlobalConfig.WriteConfig()
	if err != nil {
		log.GetLogger().Error("Modify password failed: ", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
