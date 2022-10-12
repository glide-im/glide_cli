package cli

import (
	"github.com/glide-im/glide-gui/apis"
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"testing"
	"time"
)

func TestGetToken(t *testing.T) {
	jwt := jwt_auth.NewAuthorizeImpl("secret")
	token, err := jwt.GetToken(&jwt_auth.JwtAuthInfo{
		UID:         "2",
		Device:      "1",
		ExpiredHour: 10,
	})
	if err != nil {
		panic(err)
	}
	t.Log(token.Token)
}

func TestRun(t *testing.T) {
	apis.SetBaseUrl("https://intercom.ink/api/")
	cli := NewClient("ws://localhost:8080/ws")
	err := cli.LoginByPassword("dengzii@foxmail.com", "password")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(3)
}

func TestWebsocket(t *testing.T) {

	cli := NewClient("ws://localhost:8080/ws")
	err := cli.LoginByToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjU1ODAzOTEsInVpZCI6IjEiLCJkZXZpY2UiOiIxIiwidmVyIjoxNjY1NTQ0MzkxfQ.YSWIPEN9Y9wIN3leFx8jQas-xkn-5y4mSDUgsA_uavI")
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 5)
	_, err = cli.SendTextMessage("chan1", true, "hello")
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 13)
}

func TestWebsocket2(t *testing.T) {

	cli := NewClient("ws://localhost:8080/ws")
	err := cli.LoginByToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjU1ODA0MjYsInVpZCI6IjIiLCJkZXZpY2UiOiIxIiwidmVyIjoxNjY1NTQ0NDI2fQ.aDJsYTm11QUwDBEHIqgm9JuaY9epsKBJAFGWlYRqi80")
	if err != nil {
		t.Error(err)
	}

	//time.Sleep(time.Second * 5)
	//_, err = cli.SendTextMessage("chan1", true, "hello")
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 13)
}
