package common

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/gorilla/websocket"
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

func (p *Proto) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
func (p *Proto) FromBytes(b []byte) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(p)
}
