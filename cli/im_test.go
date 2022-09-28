package cli

import (
	"github.com/glide-im/glide-gui/apis"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	apis.SetBaseUrl("https://intercom.ink/api/")
	cli := NewClient("ws://localhost:8080/ws")
	err := cli.LoginByPassword("dengzii@foxmail.com", "password")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(3)
}
