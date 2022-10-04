package api

import (
	"chat/pkg/e"
	"chat/service"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func UserRegister(c *gin.Context) {
	var userRegisterService service.UserRegisterService

	// 数据绑定
	if err := c.ShouldBind(&userRegisterService); err != nil {
		res := userRegisterService.Register()
		c.JSON(e.SUCCESS, res)
	} else {
		c.JSON(e.InvalidParams, ErrorResponse(err))
		logging.Info(err)
	}

}
