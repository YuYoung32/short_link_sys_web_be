/**
 * Created by YuYoung on 2023/3/23
 * Description: 快捷成功响应
 */

package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type MsgResponse struct {
	Msg string `json:"msg"`
}

func SuccessGeneralResp(ctx *gin.Context) {
	msg := "ok"
	ctx.JSON(http.StatusBadRequest, MsgResponse{
		Msg: msg,
	})
}
