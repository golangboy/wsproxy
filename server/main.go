package main

import (
	"encoding/base64"
	"flag"
	"log"
	"main/common"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	listenAddr := flag.String("listen", ":80", "listen address")
	flag.Parse()
	g := gin.Default()
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	g.GET("/chat", func(context *gin.Context) {
		upgrade := websocket.Upgrader{}
		conn, err := upgrade.Upgrade(context.Writer, context.Request, nil)
		ws := common.WsConn{Conn: conn}
		if err != nil {
			return
		}
		go func() {
			msgId := ""
			p := common.Proto{}
			var target net.Conn
			for {
				err = ws.Conn.ReadJSON(&p)
				if err != nil {
					break
				}
				if p.MsgType == common.ReqConnect {
					msgId = p.MsgId
					bytes, _ := base64.StdEncoding.DecodeString(string(p.Data))
					target, err = net.Dial("tcp", string(bytes))
					if err != nil {
						log.Printf("[%s]Failed to connect to the target server: %v", msgId, err)
						break
					}
					err = ws.Conn.WriteJSON(common.Proto{
						MsgType: common.ReqConnect,
						MsgId:   msgId,
					})
					if err != nil {
						log.Printf("[%s]Failed to send connection success message to the client: %v", msgId, err)
						break
					}
					go func() {
						for {
							buf := make([]byte, 10240)
							n, err := target.Read(buf[:])
							if err != nil {
								log.Printf("[%s]Failed to read data from the target server: %v", msgId, err)
								break
							}
							ws.Lock()
							err = ws.Conn.WriteJSON(common.Proto{
								MsgType: common.ReqData,
								Data:    buf[:n],
								MsgId:   msgId,
							})
							ws.Unlock()
							if err != nil {
								log.Printf("[%s]Failed to send data to the client: %v", msgId, err)
								break
							}
						}
						target.Close()
						log.Printf("[%s]Connection closed_2", msgId)
					}()
				} else if p.MsgType == common.ReqData {
					_, err = target.Write(p.Data)
					if err != nil {
						log.Printf("[%s]Failed to send data to the target server：%v", msgId, err)
						break
					}
				}
			}
			ws.Lock()
			ws.Conn.Close()
			ws.Unlock()
			log.Printf("[%s]Connection closed_1", msgId)
		}()
		//context.JSON(200, gin.H{
		//	"message": "ok",
		//})
	})
	g.Run(*listenAddr)
}
