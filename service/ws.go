package service

import (
	"chat/cache"
	"chat/conf"
	"chat/pkg/e"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

const month = 60 * 60 * 24 * 30

// 发送消息的类型
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// 回复的消息
type ReplyMsg struct {
	From    string `json:"form"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// 用户类
type Client struct {
	ID     string
	SendId string
	Socket *websocket.Conn
	Send   chan []byte
}

// 广播列，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

// Message 信息转json 包括（发送者、接收者、内容）
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *Broadcast),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateId(uid, toUid string) string {
	return uid + "->" + toUid // 1=>2
}

func Handler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")
	conn, err := (&websocket.Upgrader{

		// 检查解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil) // 链接升级为ws协议

	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// 创建一个用户实例
	client := &Client{
		ID:     CreateId(uid, toUid),
		SendId: CreateId(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}

	// 用户注册到用户管理上
	Manager.Register <- client
	go client.Read()
	go client.Write()

}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for true {
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)
		// 获取消息
		err := c.Socket.ReadJSON(&sendMsg)
		if err != nil {
			fmt.Println("数据格式不正确", err)
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}

		if sendMsg.Type == 1 {
			// 发送消息
			r1, _ := cache.RedisClient.Get(c.ID).Result() // 1->2
			r2, _ := cache.RedisClient.Get(c.SendId).Result()
			if r1 > "3" && r2 == "" {
				// 1 给2 发消息，发了三条，但是2没有回，或者没有看到，就停止发送
				replyMsg := ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "达到限制",
				}

				// json 序列化
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				cache.RedisClient.Incr(c.ID)
				_, _ = cache.RedisClient.Expire(c.ID, time.Hour*2*30*3).Result()
				// 防止分手过快，简历三个月链接
			}
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == 2 {
			// 获取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content) // string  to int
			if err != nil {
				timeT = 9999999
			}
			results, _ := FindMany(conf.MongoDBName, c.SendId, c.ID, int64(timeT), 10)
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
		}
	}

}
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for true {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccessMessage,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
