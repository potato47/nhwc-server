package server

import (
	"github.com/labstack/echo"
	"fmt"
	"encoding/json"
	"golang.org/x/net/websocket"
	"nhwc-server/data"
)

var (
	clientsMap map[int]*Client
	globalCid  int = 0
)

func init() {
	clientsMap = make(map[int]*Client)
}

func OnClientConnect(c echo.Context) error {
	globalCid++
	websocket.Handler(func(ws *websocket.Conn) {
		cid := globalCid
		cl := onConnectSuccess(cid, ws)
		fmt.Println("客户端连接当前连接客户端数:", len(clientsMap))
		defer onClientDisconnct(cl)
		for {
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				//log.Fatal(err)
				fmt.Println(err)
				break
			}
			fmt.Println("接受消息", msg)
			fmt.Println("客户端列表-------------------------------")
			myLog, _ := json.Marshal(clientsMap)
			fmt.Println(string(myLog))
			fmt.Println("用户列表-------------------------------")
			myLog, _ = json.Marshal(data.UserMap)
			fmt.Println(string(myLog))
			fmt.Println("房间列表-------------------------------")
			myLog, _ = json.Marshal(data.RoomMap)
			fmt.Println(string(myLog))
			var m WsMessage
			json.Unmarshal([]byte(msg), &m)
			// router
			switch m.Name {
			case "login":
				fmt.Println("客户端请求登录")
				onLogin(cid, m.Value)
			case "hall":
				fmt.Println("客户端请求大厅信息")
				onHall(cid, m.Value)
			case "create":
				fmt.Println("客户端请求创建房间")
				onCreate(cid, m.Value)
			case "enter":
				fmt.Println("客户端请求进入房间")
				onEnter(cid, m.Value)
			case "room":
				fmt.Println("客户端请求房间信息")
				onRoom(cid, m.Value)
			case "message":
				fmt.Println("客户端发送信息")
				onMessage(cid, m.Value)
			case "draw":
				fmt.Println("客户端同步绘画数据")
				onDraw(cid, m.Value)
			case "ready":
				fmt.Println("客户端准备")
				onReady(cid, m.Value)
			case "answer":
				fmt.Println("客户端答题")
				onAnswer(cid, m.Value)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func Bind(cid int, uid string) {
	c, ok := clientsMap[cid]
	if !ok {
		fmt.Println("没有这个client")
	} else {
		fmt.Println(uid, "号用户与", cid, "号客户端绑定")
		c.Uid = uid;
	}
}

func Send(cid int, data WsMessage) {
	msg, err := json.Marshal(data)
	fmt.Println(string(msg))
	if err != nil {
		fmt.Println(err)
	} else {
		pClient, ok := clientsMap[cid]
		if ok {
			ws := pClient.Conn
			websocket.Message.Send(ws, string(msg))
		}
	}
}

func BroadCast(roomId int, msg WsMessage, exclude int) {
	fmt.Println("广播消息", msg)
	if roomId == 0 {
		for _, user := range data.UserMap {
			Send(user.Cid, msg);
		}
	} else {
		pRoom,_ := data.GetPRoom(roomId)
		for _, user := range pRoom.SeatMap {
			if user != nil && user.Cid != exclude {
				Send(user.Cid, msg)
			}
		}
	}
}
