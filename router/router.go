package router

import (
	"chat/api"
	"chat/pkg/e"
	"chat/service"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery(), gin.Logger())
	// Recover 中间件会回复（recovers）任何恐慌（panics）如果存在恐慌，中间件会写入500，
	//	中间件还是很必要的，因为当你程序某些异常情况没有考虑到的时候，程序就退出了，服务就停止了
	// Logger 日志

	v1 := r.Group("/")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(e.SUCCESS, "success")
		})
		v1.POST("user/register", api.UserRegister)
		v1.GET("ws", service.Handler)
	}

	return r
}
