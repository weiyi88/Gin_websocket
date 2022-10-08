package ws

type Trainer struct {
	Content   string `bson:"content"`   // 内容
	StartTIme int64  `bson:"startTIme"` // 创建时间
	EndTIme   int64  `bson:"endTIme"`   // 过期时间
	Read      uint   `bson:"read"`      // 已读

}

type Reuslt struct {
	StartTime int64
	Msg       string
	Content   any
	Form      string
}