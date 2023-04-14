/**
 * Created by YuYoung on 2023/3/23
 * Description: 快捷错误响应和记录
 */

package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short_link_sys_web_be/log"
)

func ErrMissArgsResp(ctx *gin.Context) {
	msg := "miss args"
	log.GetLoggerWithSkip(2).Debug(msg)
	ctx.JSON(http.StatusBadRequest, MsgResponse{
		Msg: msg,
	})
}

func ErrInvalidArgsResp(ctx *gin.Context) {
	msg := "invalid args: out of range or invalid type"
	log.GetLoggerWithSkip(2).Debug(msg)
	ctx.JSON(http.StatusBadRequest, MsgResponse{
		Msg: msg,
	})
}

func ErrInternalResp(ctx *gin.Context) {
	msg := "internal error"
	log.GetLoggerWithSkip(2).Debug(msg)
	ctx.JSON(http.StatusInternalServerError, MsgResponse{
		Msg: msg,
	})
}

func ErrNoAuthResp(ctx *gin.Context) {
	msg := "no auth"
	log.GetLoggerWithSkip(2).Debug(msg)
	ctx.JSON(http.StatusUnauthorized, MsgResponse{
		Msg: msg,
	})
}

func ErrLoginFailedResp(ctx *gin.Context) {
	msg := "login failed"
	log.GetLoggerWithSkip(2).Debug(msg)
	ctx.JSON(http.StatusUnauthorized, MsgResponse{
		Msg: msg,
	})
}
