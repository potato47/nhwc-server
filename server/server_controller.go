package server

import (
	"time"
	"nhwc-server/data"
	"nhwc-server/config"
)

func CreateRoom(seatSum int, gameSum int) *data.Room {
	seatMap := make(map[int]*data.User)
	rid := data.GetANewRid()
	for i := 1; i < config.SeatSum; i++ {
		seatMap[i] = nil
	}
	r := data.Room{
		Rid:        rid,
		SeatMap:    seatMap,
		Painter:    0,
		WordIndex:  0,
		State:      data.RoomState_None,
		GameNum:    0,
		GameTime:   config.GameTime,
		ResultTime: config.ResultTime,
	}
	data.SetPRoom(&r)
	return &r
}

func EnterRoom(pUser *data.User, pRoom *data.Room, seat int) {
	pUser.Rid = pRoom.Rid
	pUser.Seat = seat
	pRoom.SeatMap[seat] = pUser
	if pRoom.State == data.RoomState_None {
		pRoom.State = data.RoomState_Ready
	}
}

func LeaveUser(rid int, uid string) {
	pUser := data.GetPUser(uid)

	n := 0
	for _, seat := range data.RoomMap[rid].SeatMap{
		if seat != nil {
			n++
		}
	}
	if n <= 2 {
		OverGame(rid)
	} else{
		data.RoomMap[rid].SeatMap[pUser.Seat] = nil;

		BroadCast(pUser.Rid, WsMessage{
			Name:  "exit",
			Value: pUser.Seat,
		}, pUser.Cid)
		BroadCast(pUser.Rid, WsMessage{
			Name:  "message",
			Value: pUser.Username+"退出游戏",
		}, pUser.Cid)
		pUser.Rid = 0
		pUser.Score = 0
		pUser.Seat = 0
		pUser.IsReady = false
	}
}

func GuessRight(rid int, guesserSeat int) {
	pRoom := data.RoomMap[rid];
	pRoom.SeatMap[pRoom.Painter].Score += 2
	pRoom.SeatMap[guesserSeat].Score += 1
}

func CanStartGame(rid int) bool {
	seatSum := 0
	readySum := 0
	seatMap := data.RoomMap[rid].SeatMap
	for _, seat := range seatMap {
		if seat != nil {
			seatSum++
			if seat.IsReady {
				readySum++
			}
		}
	}
	if seatSum == readySum && seatSum >= 2 {
		return true
	} else {
		return false
	}
}

func StartGame(rid int) {
	pRoom, _ := data.GetPRoom(rid)
	pRoom.Painter = getNextSeat(rid)
	pRoom.WordIndex++
	pRoom.Word = data.WordsList[pRoom.WordIndex].Word
	pRoom.Hint = data.WordsList[pRoom.WordIndex].D
	pRoom.State = data.RoomState_Draw
	pRoom.GameNum++
	pRoom.StartTime = time.Now().Second()
	timer := time.NewTimer(time.Second * time.Duration(pRoom.GameTime))
	pUser := pRoom.SeatMap[pRoom.Painter]
	cid := pUser.Cid
	Send(cid, WsMessage{
		Name:  "startMe",
		Value: pRoom,
	})
	BroadCast(pUser.Rid, WsMessage{
		Name:  "startOther",
		Value: pRoom,
	}, cid)
	go func() {
		<-timer.C
		if(pRoom != nil && pRoom.State == data.RoomState_Draw) {
			showAnswer(rid)
		}
	}()
}

func OverGame(rid int) {
	pRoom, _ := data.GetPRoom(rid)
	BroadCast(rid, WsMessage{
		Name:  "over",
		Value: data.RoomMap[rid].Word,
	}, 0)
	BroadCast(rid, WsMessage{
		Name:  "message",
		Value: "【系统】全部游戏结束,答案是:" + data.RoomMap[rid].Word,
	}, 0)
	for i, user := range pRoom.SeatMap {
		user.Rid = 0
		user.IsReady = false
		user.Seat = 0
		user.Score = 0
		pRoom.SeatMap[i] = nil
	}
	pRoom.Word = ""
	pRoom.GameNum = 0
	pRoom.Painter = 0
	pRoom.WordIndex = 0
	pRoom.Hint = ""
	pRoom.State = data.RoomState_None
}

func showAnswer(rid int) {
	BroadCast(rid, WsMessage{
		Name:  "result",
		Value: data.RoomMap[rid].Word,
	}, 0)
	BroadCast(rid, WsMessage{
		Name:  "message",
		Value: "【系统】本轮游戏结束,答案是:" + data.RoomMap[rid].Word,
	}, 0)
	timer2 := time.NewTimer(time.Second * time.Duration(data.RoomMap[rid].ResultTime))
	go func() {
		<-timer2.C
		if data.RoomMap[rid] == nil || data.RoomMap[rid].State != data.RoomState_Result {
			return
		}
		if data.RoomMap[rid].GameNum >= data.RoomMap[rid].GameTime {
			OverGame(rid)
		} else {
			StartGame(rid)
		}
	}()
}

func getNextSeat(rid int) int {
	pRoom, _ := data.GetPRoom(rid)
	currSeat := pRoom.Painter
	i := currSeat
	for {
		if i == len(pRoom.SeatMap) {
			i = 1
		} else {
			i++
		}
		if pRoom.SeatMap[i] != nil {
			return i
		}
	}

}

func IsOnlie(uid string) bool {
	_, ok := data.UserMap[uid]
	return ok
}
