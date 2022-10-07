package main

import (
	"chat/conf"
	"chat/router"
	"chat/service"
)

func main() {
	// 默认配置运行
	conf.Init()

	// 监听板块
	go service.Manager.Start()

	// 路由运行
	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)

}
