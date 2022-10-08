package service

import (
	"chat/conf"
	"chat/model/ws"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

func FirsFindtMsg(database string, sendId string, id string) (results []ws.Reuslt, err error) {
	// 首次查询（把对方发来的所有未读的都取出来
	var resutleMe []ws.Trainer
	var resultYou []ws.Trainer
	sendIdCollection := conf.MongoDBClinet.Database(database).Collection(sendId)
	idCollection := conf.MongoDBClinet.Database(database).Collection(id)
	filter := bson.M{"read": bson.M{
		"&all": []uint{0},
	}}

	sendIdCursor, err := sendIdCollection.Find(context.TODO(),
		filter, options.Find().SetSort(bson.D{
			{"startTime", 1},
		}), options.Find().SetLimit(1))

	if sendIdCursor == nil {
		return
	}
	var unReads []ws.Trainer
	err = sendIdCursor.All(context.TODO(), &unReads)
	if err != nil {
		log.Println("sendIdCursor err", err)
	}
	if len(unReads) > 0 {
		timeFIlter := bson.M{
			"startTime": bson.M{
				"$gte": unReads[0].StartTIme,
			},
		}

		sendIdTimeCursor, _ := sendIdCollection.Find(context.TODO(), timeFIlter)
		idTimeCursor, _ := idCollection.Find(context.TODO(), timeFIlter)
		err = sendIdTimeCursor.All(context.TODO(), &resultYou)
		err = idTimeCursor.All(context.TODO(), &resutleMe)

		results, err = AppendAndSort(resutleMe, resultYou)

	} else {
		results, err = FindMany(database, sendId, id, 999999999, 10)
	}
	overTimeFileter := bson.D{
		{
			"$and", bson.A{
				bson.D{{"endTime", bson.M{"&lt": time.Now().Unix()}}},
				bson.D{{"read", bson.M{"$eq": 1}}},
			}},
	}

	_, _ = sendIdCollection.DeleteMany(context.TODO(), overTimeFileter)
	_, _ = idCollection.DeleteMany(context.TODO(), overTimeFileter)

	// 将所有的纬度设置为已读

	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{"read": 1},
	})

	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{"endTime": time.Now().Unix() + int64(3*month)},
	})
	return
}

// 排序
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
