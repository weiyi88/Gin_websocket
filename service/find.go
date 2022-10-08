package service

import (
	"chat/conf"
	"chat/model/ws"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

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

func FindMany(database, sendId, id string, time int64, pageSize int) (results []ws.Reuslt, err error) {
	var resultMe []ws.Trainer  // id
	var resultYou []ws.Trainer // sendID
	sendIdCollection := conf.MongoDBClinet.Database(database).Collection(sendId)
	idCollection := conf.MongoDBClinet.Database(database).Collection(id)

	// 如果不知道用什么context ，可以通过context.TODO() 产生context
	sendIDTimeCuror, err := sendIdCollection.Find(context.TODO(),
		options.Find().SetSort(bson.D{{"startTime", -1}}),
		options.Find().SetLimit(int64(pageSize)))

	idTimeCurcor, err := idCollection.Find(context.TODO(),
		options.Find().SetSort(bson.D{{"startTime", -1}}),
		options.Find().SetLimit(int64(pageSize)))

	err = sendIDTimeCuror.All(context.TODO(), &resultYou)
	err = idTimeCurcor.All(context.TODO(), &resultMe)

	results, _ = AppendAndSort(resultMe, resultYou)
	return
}

func AppendAndSort(resultMe, resultYou []ws.Trainer) (results []ws.Reuslt, err error) {
	for _, r := range resultMe {

		// 返回消息的结构体
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTIme,
		}

		// 作为返回结果
		result := ws.Reuslt{
			StartTime: r.StartTIme,
			Msg:       fmt.Sprintf("%v", sendSort),
			Form:      "me",
		}

		results = append(results, result)
	}
	return
}
