package common

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Proto struct {
	MsgType int    `json:"msg_type"`
	MsgId   string `json:"msg_id"`
	Data    []byte `json:"data"`
}
type WsConn struct {
	Conn *websocket.Conn
	sync.Mutex
}

const (
	ReqConnect = iota
	ReqData    = iota
)
