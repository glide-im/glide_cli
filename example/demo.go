package main

import (
	"github.com/glide-im/glide-gui/cli"
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"github.com/glide-im/glide/pkg/logger"
	"time"
)

func main() {

	client := cli.NewClient("ws://localhost:8080/ws")
	err := client.LoginByToken(getToken("1"))
	if err != nil {
		panic(err)
	}

	message, err := client.SendTextMessage("2", "Hello From uid 1")
	if err != nil {
		panic(err)
	}

	logger.D("message id: %s", message.Mid)

	time.Sleep(time.Minute)
}

func getToken(uid string) string {
	/// simulate login success
	jwt := jwt_auth.NewAuthorizeImpl("your_secret_here")
	token, err := jwt.GetToken(&jwt_auth.JwtAuthInfo{
		UID:         uid,
		Device:      "1",
		ExpiredHour: 10,
	})
	if err != nil {
		panic(err)
	}
	return token.Token
}
