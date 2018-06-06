package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	UserData struct {
		Uid      string `json:"uid"`
		Username string `json:"username"`
		Password string `json:"password"`
		Admin    bool   `json:"admin"`
	}
)

func init() {

}

func ValidateUserData(uid string, password string) (*UserData, error) {
	session, err := mgo.Dial("119.29.40.244:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("users");
	user := UserData{}
	err = c.Find(bson.M{"uid": uid, "password": password}).One(&user)
	return &user, err
}

//
//func getUser(uid string, pwd string)(*UserData, error) {
//	if uid == "001" && pwd == "001" {
//		return &UserData{
//			Username: "next",
//			Password: "001",
//			Uid:"001",
//			Admin:true,
//		},errors.New("用户名或密码错误")
//	}
//}
