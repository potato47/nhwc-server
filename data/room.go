package data

import (
	"fmt"
	"nhwc-server/config"
)

const (
	RoomState_None   = 0
	RoomState_Ready  = 1
	RoomState_Draw   = 2
	RoomState_Result = 3
	RoomState_Over   = 4
)

type Room struct {
	Rid        int           `json:"rid"`
	SeatMap    map[int]*User `json:"seatMap"`
	Painter    int           `json:"painter"`
	Word       string        `json:"word"`
	WordIndex  int           `json:"wordIndex"`
	Hint       string        `json:"hint"`
	StartTime  int           `json:"startTime"`
	State      int           `json:"state"`
	GameNum    int           `json:"gameNum"`
	GameTime   int           `json:"gameTime"`
	ResultTime int           `json:"resultTime"`
	//Num             int
	//Status          int
	//PainterFd       int
	//LastReadyTime   int
	//ResultStartTime int
	//Word            string
	//Hint            string
	//WordIndex       int
	//Fd              map[int]int
}

var RoomMap map[int]*Room

func init() {
	RoomMap = make(map[int]*Room)
	for i := 1; i <= config.RoomSum; i++ {
		seatMap := make(map[int]*User)
		for j := 1; j <= config.SeatSum; j++ {
			seatMap[j] = nil
		}
		RoomMap[i] = &Room{
			Rid:       i,
			SeatMap:   seatMap,
			Painter:   0,
			WordIndex: 0,
			State:     RoomState_None,
			GameNum:   0,
			GameTime: config.GameTime,
			ResultTime: config.ResultTime,
		}
	}
}

func SetPRoom(pRoom *Room) {
	RoomMap[pRoom.Rid] = pRoom
}

func GetPRoom(rid int) (*Room,bool) {
	pRoom, ok := RoomMap[rid]
	if !ok {
		fmt.Println("没有该房间")
	}
	return pRoom, ok
}

func GetRoomIdList() []interface{} {
	var list []interface{}
	for _, r := range RoomMap {
		l := 0;
		for _, s := range r.SeatMap {
			if s != nil {
				l++;
			}
		}
		list = append(list, struct {
			Rid int `json:"rid"`
			Num int `json:"num"` //人数
			Max int `json:"max"` //房间座位数
		}{
			Rid: r.Rid,
			Num: l,
			Max: len(r.SeatMap),
		})
	}
	fmt.Println("房间数据")
	fmt.Println(list)
	return list
}

func GetFreeSeat(pRoom *Room) int {
	seatMap := pRoom.SeatMap
	for i, _ := range seatMap {
		if seatMap[i] == nil {
			return i
		}
	}
	return 0
}

func GetANewRid() int {
	maxRid := 0
	for _,v := range RoomMap {
		if v.Rid > maxRid {
			maxRid = v.Rid
		}
	}
	maxRid++
	return maxRid
}
