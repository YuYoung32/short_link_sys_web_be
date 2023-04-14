/**
 * Created by YuYoung on 2023/4/7
 * Description: 服务器本身信息
 */

package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func StaticInfoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, staticInfo)
}
