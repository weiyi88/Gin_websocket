package service

import (
	"chat/conf"
	"chat/model/ws "
	"context"
	"time"
)

func InsertMsg(database, id, content string, read uint, expire int64) error {
	// 插入 mongodb中
	collection := conf.MongoDBClinet.Database(database).Collection(id)

	// 插入数据格式规范
	comment := ws.Trainer{
		Content:   content,
		StartTIme: time.Now().Unix(),
		EndTIme:   time.Now().Unix() + expire,
		Read:      read,
	}

	_, err := collection.InsertOne(context.TODO(), comment)

	return err

}
