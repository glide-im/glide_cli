package main

import (
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
)

type Client struct {
	ws GlideWsClient
}

func NewClient(url string) *Client {
	ws, err := NewGlideWsClient(url)
	if err != nil {
		panic(err)
	}
	c := &Client{
		ws: ws,
	}
	c.init()
	go ws.Run()
	return c
}

func (c *Client) init() {
	c.ws.ListenerMessage(func(m *messages.GlideMessage) {
		switch m.GetAction() {
		case actionChatMessage:
			cm := messages.ChatMessage{}
			err := m.Data.Deserialize(&cm)
			if err != nil {
				logger.ErrE("deserialize msg", err)
			}
			err = c.ws.Send(messages.NewMessage(0, actionAckRequest, &messages.AckRequest{
				Mid:  cm.Mid,
				From: cm.From,
			}))
			if err != nil {
				logger.ErrE("ack request", err)
			}
			logger.D("chat message: %v", cm)
		case actionAckNotify:

		}
	})
}

func (c *Client) LoginByPassword(account, password string) error {
	return nil
}

func (c *Client) LoginByToken(token string) error {
	return c.ws.Auth(token)
}
