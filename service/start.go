package service

import (
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
		}
	}
}