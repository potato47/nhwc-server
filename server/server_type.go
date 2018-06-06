package server

import "golang.org/x/net/websocket"

type (
	Client struct {
		Cid  int
		Uid  string
		Conn *websocket.Conn
	}
	Token struct {
		Token string `json:"token"`
	}
	WsMessage struct {
		Name  string `json:"name"`
		Value interface{} `json:"value"`
	}
	Response struct {
		Data  interface{} `json:"data"`
		Error interface{} `json:"error"`
	}
	UserData struct {
		Uid      string `json:"uid"`
		Username string `json:"username"`
		Token    string `json:"token"`
	}
	ErrorData struct {
		ErrorCode    int    `json:"errorcode"`
		ErrorMessage string `json:"errormessage"`
	}
)
