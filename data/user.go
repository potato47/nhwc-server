package data

import "fmt"

type User struct {
	Uid      string `json:"uid"`
	Cid      int    `-`
	Rid      int    `json:"rid"`
	Seat     int    `json:"seat"`
	Score    int    `json:"score"`
	IsReady  bool   `json:"isReady"`
	Username string `json:"username"`
}

var UserMap map[string]*User

func init() {
	UserMap = make(map[string]*User)
}

func SetPUser(pUser *User) {
	UserMap[pUser.Uid] = pUser
}

func GetPUser(uid string) *User {
	pUser, ok := UserMap[uid]
	if !ok {
		fmt.Println("没有该用户")
	}
	return pUser
}
