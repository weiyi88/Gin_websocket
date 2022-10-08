package service

import (
	"chat/conf"
	"chat/pkg/e"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

func (manage *ClientManager) Start() {
	// 监听管道通信
	for true {
		fmt.Println("------监听管道通信------")
		select {
		case conn := <-Manager.Register:
			fmt.Printf("有新的链接 &v", conn.ID)

			// 链接放到用户管理上
			Manager.Clients[conn.ID] = conn
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已经链接到服务器",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)

		// 链接中断
		case conn := <-Manager.Unregister:
			fmt.Printf("链接失败 %s", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "链接中断",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}

		// 广播给用户
		case broadcast := <-Manager.Broadcast: // 1 -> 2
			message := broadcast.Message
			sendId := broadcast.Client.SendId // 2 ->1
			flag := false                     // 默认对方不在线

			// 如果不在线
			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}

			// 在线状态
			id := broadcast.Client.ID // 1-> 2
			if flag {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答中",
				}
				msg, err := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				// 消息插入到mongodb中
				err = InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month)) // 1代表已读
				if err != nil {
					fmt.Println("Insert One err ", err)
				}
			}

		}
	}
}
