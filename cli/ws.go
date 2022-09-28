package cli

import (
	"encoding/json"
	"errors"
	"github.com/glide-im/glide/pkg/auth"
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
)

const (
	actionChatMessage = "message.chat"
	actionHeartbeat   = "heartbeat"

	actionApiAuth     = "api.auth"
	actionNotifyError = "notify.error"
	actionAckMessage  = "ack.message"
	actionAckRequest  = "ack.request"
	actionAckNotify   = "ack.notify"
)

type GlideWsClient interface {

	// Run 开始接受处理消息
	Run() error

	// Close 关闭客户端, 断开连接
	Close() error

	// Auth  token 登录
	Auth(token string) error

	// ListenerMessage 监听消息
	ListenerMessage(l func(m *messages.GlideMessage))

	// SendChatMessage 发送聊天消息
	// 等待服务器确认收到后返回带服务器返回消息id, 失败时返回错误
	SendChatMessage(to string, typ int32, content string) (*messages.ChatMessage, error)

	// SendApiMessage 发送 Api 消息
	// 返回结果, 失败返回错误
	SendApiMessage(action string, data interface{}) (*messages.Data, error)

	Send(i interface{}) error
}

type glide struct {
	uid   string
	token string

	seq      int64
	conn     *websocket.Conn
	messages chan *messages.GlideMessage

	hbTicker *time.Ticker

	listener map[interface{}]interface{}
}

func NewGlideWsClient(url string) (GlideWsClient, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout:  3 * time.Second,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
		Jar:               nil,
	}

	dial, response, err := dialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 101 {
		return nil, errors.New("failed to connect to server, " + strconv.Itoa(response.StatusCode) + " " + response.Status)
	}
	g := &glide{
		seq:      100,
		conn:     dial,
		messages: make(chan *messages.GlideMessage, 10),
		listener: map[interface{}]interface{}{},
	}
	return g, nil
}

func (g *glide) Run() error {
	go g.handleMessage()
	go g.startHeartbeat()
	g.recv()
	return nil
}

func (g *glide) Close() error {
	return g.conn.Close()
}

func (g *glide) Auth(token string) error {

	resp, err := g.SendApiMessage(actionApiAuth, &auth.Token{Token: token})
	if err != nil {
		return err
	}
	data := &jwt_auth.Response{}
	if err = resp.Deserialize(data); err != nil {
		return err
	}
	g.uid = data.Uid
	g.token = data.Token
	return nil
}

func (g *glide) ListenerMessage(l func(m *messages.GlideMessage)) {
	g.listener[time.Now().UnixNano()] = l
}

func (g *glide) Send(i interface{}) error {
	return g.write(i)
}

func (g *glide) SendChatMessage(to string, typ int32, content string) (*messages.ChatMessage, error) {
	cm := &messages.ChatMessage{
		CliMid:  uuid.New().String(),
		From:    g.uid,
		To:      to,
		Type:    typ,
		Content: content,
		SendAt:  time.Now().Unix(),
	}
	err := g.send(to, actionChatMessage, cm)
	if err != nil {
		return nil, err
	}
	ack := g.waitAck(cm.CliMid)
	cm.Mid = ack.Mid
	return cm, nil
}

func (g *glide) SendApiMessage(action string, data interface{}) (*messages.Data, error) {
	seq := g.seq
	gm := messages.GlideMessage{
		Ver:    1,
		Seq:    seq,
		Action: action,
		Data:   messages.NewData(data),
	}
	err := g.write(&gm)
	if err != nil {
		return nil, err
	}
	g.seq++
	resp := g.waitSeq(seq)
	return resp.Data, nil
}

func (g *glide) send(to string, action string, data interface{}) error {
	gm := messages.GlideMessage{
		Ver:    1,
		Seq:    g.seq,
		Action: action,
		From:   g.uid,
		To:     to,
		Data:   messages.NewData(data),
	}
	err := g.write(&gm)
	g.seq++
	return err
}

func (g *glide) waitSeq(seq int64) *messages.GlideMessage {
	ch := make(chan *messages.GlideMessage)
	defer close(ch)

	l := func(message *messages.GlideMessage) {
		if message.GetSeq() == seq {
			ch <- message
		}
	}
	g.listener[seq] = l
	res := <-ch
	delete(g.listener, seq)
	return res
}

func (g *glide) waitAck(cliId string) *messages.AckMessage {
	ch := make(chan *messages.AckMessage)
	defer close(ch)

	l := func(message *messages.GlideMessage) {
		if message.Action == actionAckMessage {
			ack := messages.AckMessage{}
			err := message.Data.Deserialize(&ack)
			if err != nil {
				return
			}
			if ack.CliMid == cliId {
				ch <- &ack
			}
		}
	}
	g.listener[cliId] = l
	res := <-ch
	delete(g.listener, cliId)

	return res
}

func (g *glide) startHeartbeat() {
	g.hbTicker = time.NewTicker(time.Second * 30)
	for range g.hbTicker.C {
		err := g.write(messages.NewMessage(g.seq, actionHeartbeat, nil))
		if err != nil {
			logger.E("heartbeat failed: %v", err)
		}
	}
}

func (g *glide) handleMessage() {
	for m := range g.messages {
		for _, l := range g.listener {
			go func(li func(*messages.GlideMessage), msg *messages.GlideMessage) {
				defer func() {
					if r := recover(); r != nil {
						logger.E("%v", r)
					}
				}()
				li(msg)
			}(l.(func(*messages.GlideMessage)), m)
		}
	}
	logger.D("receive stopped")
}

func (g *glide) write(i interface{}) error {
	bytes, err := json.Marshal(i)
	if err != nil {
		return err
	}
	logger.D("[send] %s", string(bytes))
	return g.conn.WriteMessage(websocket.TextMessage, bytes)
}

func (g *glide) recv() {
	defer func() {
		a := recover()
		if a != nil {
			logger.E("%v", a)
		}
	}()

	g.messages = make(chan *messages.GlideMessage, 10)
	errs := 0
	for {
		messageType, bytes, err := g.conn.ReadMessage()
		if messageType != websocket.TextMessage {
			continue
		}
		if err != nil {
			logger.E("received message error: %v", err)
			errs++
			if errs > 10 {
				break
			}
			continue
		}
		var message messages.GlideMessage
		err = json.Unmarshal(bytes, &message)
		if err != nil {
			logger.E("unmarshal message error: %v", err)
			continue
		}
		errs = 0
		logger.D("[recv] %s", string(bytes))
		g.messages <- &message
	}
}
