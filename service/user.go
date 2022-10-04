package service

import (
	"chat/model"
	"chat/pkg/e"
	"chat/serializer"
)

type UserRegisterService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (service *UserRegisterService) Register() serializer.Response {
	var user model.User
	count := 0
	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).First(&user).Count(&count)
	if count != 0 {
		return serializer.Response{
			Status: e.InvalidParams,
			Msg:    "用户名已存在",
		}
	}

	user = model.User{
		UserName: service.UserName,
	}

	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Status: e.ERROR,
			Msg:    "加密出错",
		}
	}

	// 创建成功
	model.DB.Create(&user)
	return serializer.Response{
		Status: e.SUCCESS,
		Msg:    "创建成功",
	}

}
