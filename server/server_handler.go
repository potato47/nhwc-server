package server

import (
	"fmt"
	"golang.org/x/net/websocket"
	"nhwc-server/data"
	"nhwc-server/model"
)

//客户端连接
func onConnectSuccess(cid int, pWs *websocket.Conn) *Client {
	// 保存所有客户端连接
	c := Client{
		Cid:  cid,
		Conn: pWs,
	}
	clientsMap[cid] = &c
	return &c
}

//客户端断开连接
func onClientDisconnct(c *Client) {
	pUser := data.GetPUser(c.Uid)
	if pUser != nil {
		if pUser.Rid > 0 {
			LeaveUser(pUser.Rid, pUser.Uid)
		}
		delete(data.UserMap, c.Uid)
		// 更新大厅房间信息
		roomList := data.GetRoomIdList()
		BroadCast(0, WsMessage{
			Name: "hall",
			Value: struct {
				RoomList []interface{} `json:"roomlist"`
			}{
				RoomList: roomList,
			},
		}, c.Cid)
	}
	delete(clientsMap, c.Cid)
	fmt.Println("客户端断开连接，当前客户端数", len(clientsMap))
	c.Conn.Close()
}

// 登录
func onLogin(cid int, msg interface{}) {
	// 验证用户登录
	parseMsg := msg.(map[string]interface{})
	username := parseMsg["username"].(string)
	password := parseMsg["password"].(string)
	pUserData, err := model.ValidateUserData(username, password)
	if err != nil {
		Send(cid, WsMessage{
			Name: "error",
			Value: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "1",
				Message: "用户名或密码错误",
			},
		})
	} else if IsOnlie(pUserData.Uid) {
		Send(cid, WsMessage{
			Name: "error",
			Value: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "2",
				Message: "该用户已经在游戏中",
			},
		})
	} else {
		u := data.User{
			Uid:      pUserData.Uid,
			Cid:      cid,
			Rid:      0,
			Seat:     0,
			Username: pUserData.Username,
		}
		data.SetPUser(&u)
		// 用户与client绑定
		Bind(cid, pUserData.Uid)

		m := WsMessage{
			Name: "login",
			Value: struct {
				Username string `json:"username"`
				Uid      string `json:"uid"`
			}{
				Username: pUserData.Username,
				Uid:      pUserData.Uid,
			},
		}
		Send(cid, m)
	}
}

// 获取大厅信息
func onHall(cid int, msg interface{}) {
	roomList := data.GetRoomIdList()
	m := WsMessage{
		Name: "hall",
		Value: struct {
			RoomList []interface{} `json:"roomlist"`
		}{
			RoomList: roomList,
		},
	}
	Send(cid, m)
}

// 创建房间
func onCreate(cid int, msg interface{}) {
	parseMsg := msg.(map[string]interface{})

	pRoom := CreateRoom(int(parseMsg["seatSum"].(float64)), int(parseMsg["gameSum"].(float64)))
	m := WsMessage{
		Name:  "create",
		Value: pRoom,
	}
	Send(cid, m)
}

// 进入房间
func onEnter(cid int, msg interface{}) {
	parseMsg := msg.(map[string]interface{})
	rid := int(parseMsg["rid"].(float64))
	uid := clientsMap[cid].Uid
	pRoom, ok := data.GetPRoom(rid)
	pUser := data.GetPUser(uid)
	if ok {
		freeSeat := data.GetFreeSeat(pRoom)
		if freeSeat == 0 {
			m := WsMessage{
				Name:  "error",
				Value: "房间已满",
			}
			Send(cid, m)
		} else {
			EnterRoom(pUser, pRoom, freeSeat)
			// 通知客户端进入房间
			m := WsMessage{
				Name:  "enter",
				Value: pRoom,
			}
			Send(cid, m)

			//通知其他客户端有人加入房间
			pUser := data.GetPUser(uid)
			BroadCast(rid, WsMessage{
				Name:  "message",
				Value: "【系统】" + pUser.Username + "进入房间",
			}, cid)
			BroadCast(rid, WsMessage{
				Name:  "room",
				Value: pRoom,
			}, cid)

			// 更新大厅房间信息
			roomList := data.GetRoomIdList()
			BroadCast(0, WsMessage{
				Name: "hall",
				Value: struct {
					RoomList []interface{} `json:"roomlist"`
				}{
					RoomList: roomList,
				},
			}, cid)
		}
	} else {
		m := WsMessage{
			Name:  "error",
			Value: "没有该房间",
		}
		Send(cid, m)
	}

}

// 请求房间信息
func onRoom(cid int, msg interface{}) {
	parseMsg := msg.(map[string]interface{})
	rid := int(parseMsg["rid"].(float64))
	pRoom, _ := data.GetPRoom(rid)
	m := WsMessage{
		Name:  "room",
		Value: pRoom,
	}
	Send(cid, m)
}

// 客户端发送信息
func onMessage(cid int, msg interface{}) {
	pUser := data.GetPUser(clientsMap[cid].Uid)
	m := WsMessage{
		Name:  "message",
		Value: "【" + pUser.Username + "】 " + msg.(string),
	}
	BroadCast(pUser.Rid, m, 0)
}

// 同步客户端绘画信息
func onDraw(cid int, msg interface{}) {
	pUser := data.GetPUser(clientsMap[cid].Uid)
	m := WsMessage{
		Name:  "draw",
		Value: msg,
	}
	BroadCast(pUser.Rid, m, cid)
}

// 客户端请求准备
func onReady(cid int, msg interface{}) {
	pUser := data.GetPUser(clientsMap[cid].Uid)
	if pUser.IsReady {
		fmt.Println("用户已经准备");
		Send(cid, WsMessage{
			Name: "error",
			Value: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "2",
				Message: "你已经准备了",
			},
		})
	} else {
		pUser.IsReady = true;
		m := WsMessage{
			Name:  "ready",
			Value: map[string]int{"seat": pUser.Seat},
		}
		BroadCast(pUser.Rid, m, 0)
		if CanStartGame(pUser.Rid) {
			StartGame(pUser.Rid)
		}
	}
}

// 客户端答题
func onAnswer(cid int, msg interface{}) {
	pUser := data.GetPUser(clientsMap[cid].Uid)
	pRoom, _ := data.GetPRoom(pUser.Rid)
	if pRoom.State == data.RoomState_Draw && msg.(string) == pRoom.Word {
		m := WsMessage{
			Name:  "message",
			Value: "【系统】 " + pUser.Username + "猜对了答案",
		}
		BroadCast(pUser.Rid, m, 0)
		//更新分数
		GuessRight(pRoom.Rid, pUser.Seat)
		//m2 := WsMessage{
		//	Name:  "score",
		//	Value: map[int]int{pUser.Seat: pUser.Score, pRoom.Painter: pRoom.SeatMap[pRoom.Painter].Score},
		//}
		//BroadCast(pUser.Rid, m2, 0)

		m3 := WsMessage{
			Name: "answer",
			Value: map[string]interface{}{
				"seat":    pUser.Seat,
				"isRight": true,
				"scores":  map[int]int{pUser.Seat: pUser.Score, pRoom.Painter: pRoom.SeatMap[pRoom.Painter].Score},
			},
		}
		BroadCast(pUser.Rid, m3, 0)
	} else {
		m := WsMessage{
			Name: "answer",
			Value: map[string]interface{}{
				"seat":    pUser.Seat,
				"isRight": false,
				"score":   pUser.Score,
			},
		}
		BroadCast(pUser.Rid, m, 0)
		m2 := WsMessage{
			Name:  "message",
			Value: "【系统】 " + pUser.Username + "猜" + msg.(string) + ",答案错误",
		}
		BroadCast(pUser.Rid, m2, 0)
	}
}

// 客户端退出
func onLeave(cid int, msg interface{}) {
	pUser := data.GetPUser(clientsMap[cid].Uid)
	LeaveUser(pUser.Rid, pUser.Uid)
}
