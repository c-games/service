package service

import (
	"testing"
	"time"
)

var (
	ws *Websocket
)
func TestNewWebSocket(t *testing.T) {

	p:=&WSParameter{
		"127.0.0.1",
		"9000",
		500,
		500,
		time.Second*time.Duration(20),
		time.Second*time.Duration(20),
		10240,
		10240,
	}
	so, err := NewWebsocket(p)

	if err != nil {
		t.Error(err)
	}else{
		ws=so
	}


}

func TestStart(t *testing.T){
	if ws!=nil {

		err := ws.Start()
		if err != nil {
			t.Error(err)
		}
	}else{
		t.Error("ws is nil")
	}
}
