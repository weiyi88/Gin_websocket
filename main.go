package main

import (
	"chat/conf"
	"chat/router"
)

func main() {
	// 默认配置运行
	conf.Init()

	// 路由运行
	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)

}
