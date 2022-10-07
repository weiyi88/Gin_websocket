package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	logging "github.com/sirupsen/logrus"
	"time"
)

var DB *gorm.DB

func Database(connString string) {
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		fmt.Println("connect err:", err)
	}
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	logging.Info("Mysql Connnect Successfully")

	if gin.Mode() == "release" {
		db.LogMode(false)
	}

	db.SingularTable(true)       // 默认不加复数s
	db.DB().SetMaxIdleConns(20)  // 设置链接池，空闲
	db.DB().SetMaxOpenConns(100) //设置打开最大链接
	db.DB().SetConnMaxLifetime(time.Second * 30)
	DB = db
	fmt.Println("herrrr")
	migration()
	fmt.Println("herrr324234r")

}
