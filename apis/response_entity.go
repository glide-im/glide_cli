package apis

import "encoding/json"

type CommonResponse struct {
	Code int
	Msg  string
	Data json.RawMessage
}

type AuthResponse struct {
	Token    string   `json:"token"`
	Uid      int64    `json:"uid"`
	Servers  []string `json:"servers"`
	NickName string   `json:"nick_name"`
	//App      app.App  `json:"app"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Device int64  `json:"device"`
}

type GuestAuthResponse struct {
	Token    string   `json:"token"`
	Uid      int64    `json:"uid"`
	Servers  []string `json:"servers"`
	AppID    int64    `json:"app_id"`
	NickName string   `json:"nick_name"`
}
