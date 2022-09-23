package main

import (
	"github.com/glide-im/glide/pkg/auth"
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"testing"
)

const (
	tokenUid1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njc0MzI5NTYsInVpZCI6IjEiLCJkZXZpY2UiOiIwIiwidmVyIjoxNjYzODMyOTU2fQ.IbTIZYm2fpjynjzpheCd719jHvymBF8GmztI5ZcxPO0"
	tokenUid2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njc0MzI5ODEsInVpZCI6IjIiLCJkZXZpY2UiOiIwIiwidmVyIjoxNjYzODMyOTgxfQ.0A-vBwMoXcy1fJ_gQUTXQomN_Cp9NPyq8vD3ARD0iLo"
)

func TestGlide(t *testing.T) {
	g, err := NewGlideWsClient("ws://localhost:8080/ws")
	if err != nil {
		t.Error(err)
	}
	go g.Run()
	resp, err := g.SendApiMessage(actionApiAuth, &auth.Token{Token: tokenUid1})
	if err != nil {
		t.Error(err)
	}
	t.Log(">>", resp.String())

	message, err := g.SendChatMessage("2", 2, "hello")

	t.Log(message, err)
}

func TestGenToken(t *testing.T) {
	j := jwt_auth.NewAuthorizeImpl("secret")
	token, _ := j.GetToken(&jwt_auth.JwtAuthInfo{
		UID:         "2",
		Device:      "0",
		ExpiredHour: 1000,
	})
	t.Log(token.Token)
}
