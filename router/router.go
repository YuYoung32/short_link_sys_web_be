/**
 * Created by YuYoung on 2023/3/22
 * Description: 总路由
 */

package router

import (
	"github.com/gin-gonic/gin"
)

func ServerRouter(router *gin.Engine) {

}
func VisitRouter(router *gin.Engine) {

}

func LinkRouter(router *gin.Engine) {

}

func LoadAllRouter(router *gin.Engine) {
	ServerRouter(router)
	VisitRouter(router)
	LinkRouter(router)
}
