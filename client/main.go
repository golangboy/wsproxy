package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"main/common"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var wsServer *string

func main() {
	os.Setenv("http_proxy", "")
	os.Setenv("https_proxy", "")
	os.Setenv("HTTP_PROXY", "")
	os.Setenv("HTTPS_PROXY", "")
	os.Setenv("ALL_PROXY", "")

	listenAddr := flag.String("listen", ":1180", "listen socks5 address")
	wsServer = flag.String("ws", "localhost:80", "websocket server address")
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", listenAddr, err)
	}

	log.Printf("SOCKS5 proxy server is listening on %s", *listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}
func convertToMap(input string) map[string]string {
	result := make(map[string]string)

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			result[strings.ToLower(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
		}
	}

	return result
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	// Step 1: Version identification and authentication
	// Read and verify the SOCKS5 initial handshake message
	buf := make([]byte, 257)
	nbytes, err := io.ReadAtLeast(clientConn, buf, 2)
	if err != nil {
		log.Printf("Failed to read client handshake: %v", err)
		return
	}
	_ = nbytes

	var destAddr string
	var destPort string
	isSocks5 := false
	isHttps := false
	// Check SOCKS version and authentication methods
	if buf[0] == 0x05 {
		isSocks5 = true
		// Number of authentication methods supported
		numMethods := int(buf[1])
		authMethods := buf[2 : 2+numMethods]

		// Check if "no authentication" method (0x00) is supported
		noAuth := false
		for _, m := range authMethods {
			if m == 0x00 {
				noAuth = true
				break
			}
		}

		if !noAuth {
			log.Printf("No supported authentication methods")
			// Send handshake failure response to client
			clientConn.Write([]byte{0x05, 0xFF})
			return
		}

		// Send handshake response to client indicating "no authentication" method
		clientConn.Write([]byte{0x05, 0x00})

		// Step 2: Request processing
		// Read and verify the SOCKS5 request
		_, err = io.ReadAtLeast(clientConn, buf, 4)
		if err != nil {
			log.Printf("Failed to read client request: %v", err)
			return
		}

		if buf[0] != 0x05 {
			log.Printf("Unsupported SOCKS version: %v", buf[0])
			return
		}

		if buf[1] != 0x01 {
			log.Printf("Unsupported command: %v", buf[1])
			return
		}

		// Check the address type
		switch buf[3] {
		case 0x01: // IPv4 address
			ip := net.IP(buf[4 : 4+net.IPv4len])
			destAddr = ip.String()
			destPort = fmt.Sprintf("%d", int(buf[8])<<8+int(buf[9]))
		case 0x03: // Domain name
			domainLen := int(buf[4])
			domain := string(buf[5 : 5+domainLen])
			destAddr = domain
			destPort = strconv.Itoa(int(buf[5+domainLen])<<8 + int(buf[5+domainLen+1]))
		case 0x04: // IPv6 address
			ip := net.IP(buf[4 : 4+net.IPv6len])
			destAddr = ip.String()
			destPort = strconv.Itoa(int(buf[20])<<8 + int(buf[21]))
		default:
			log.Printf("Unsupported address type: %v", buf[3])
			return
		}

		// Send request response to client indicating success
		clientConn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	} else if (buf[0] == 'G' && buf[1] == 'E') || (buf[0] == 'P' && buf[1] == 'O') || (buf[0] == 'O' && buf[1] == 'P') || (buf[0] == 'P' && buf[1] == 'U') || (buf[0] == 'D' && buf[1] == 'E') || (buf[0] == 'H' && buf[1] == 'E') || (buf[0] == 'T' && buf[1] == 'R') {
		r := convertToMap(string(buf[:nbytes]))
		hosts := r["host"]
		if strings.Index(hosts, ":") == -1 {
			destAddr = hosts
			destPort = "80"
		} else {
			destAddr = strings.Split(hosts, ":")[0]
			destPort = strings.Split(hosts, ":")[1]
		}
	} else if buf[0] == 'C' && buf[1] == 'O' {
		isHttps = true
		r := convertToMap(string(buf[:nbytes]))
		hosts := r["host"]
		if strings.Index(hosts, ":") == -1 {
			destAddr = hosts
			destPort = "443"
		} else {
			destAddr = strings.Split(hosts, ":")[0]
			destPort = strings.Split(hosts, ":")[1]
		}
	}
	u := url.URL{
		Scheme: "ws",
		Host:   *wsServer,
		Path:   "/chat",
	}
	msgId := uuid.NewV4().String()
	log.Printf("[%s]Client requests proxy connection: %s:%s", msgId, destAddr, destPort)
	// connection to websocket server
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("[%s]Failed to connect to the WebSocket server: %v", msgId, err)
		clientConn.Close()
		return
	}
	ws := common.WsConn{
		Conn: conn,
	}
	ws.Lock()
	err = ws.Conn.WriteMessage(websocket.BinaryMessage, (&common.Proto{
		MsgType: common.ReqConnect,
		Data:    []byte(base64.StdEncoding.EncodeToString([]byte(destAddr + ":" + destPort))),
		MsgId:   msgId,
	}).ToBytes())
	ws.Unlock()
	if err != nil {
		log.Printf("[%s]Failed to write connect request to WebSocket server: %v", msgId, err)
		clientConn.Close()
		return
	}
	connectResp := common.Proto{}

	_, b, err := ws.Conn.ReadMessage()
	connectResp.FromBytes(b)
	if err != nil {
		log.Printf("[%s]Failed to request proxy target from WebSocket server: %v", msgId, err)
		clientConn.Close()
		return
	}

	if connectResp.MsgType != common.ReqConnect {
		log.Printf("[%s]ReqConnect failed: %v", msgId, err)
		ws.Conn.Close()
		clientConn.Close()
		return
	}
	if isHttps {
		clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	}
	if false == isSocks5 && false == isHttps {
		err = ws.Conn.WriteJSON(common.Proto{MsgType: common.ReqData, Data: buf[:nbytes], MsgId: msgId})
		if err != nil {
			log.Printf("[%s]Failed to send http data to WebSocket server: %v", msgId, err)
		}
	}
	go func() {
		resp := common.Proto{}
		for {
			_, p, err := ws.Conn.ReadMessage()
			resp.FromBytes(p)
			if err != nil {
				log.Printf("[%s]Failed to read data from WebSocket server: %v", msgId, err)
				break
			}
			_, err = clientConn.Write(resp.Data)
			if err != nil {
				log.Printf("[%s]Failed to write data to the client: %v", msgId, err)
				break
			}
		}
	}()
	for {
		var buf [1024 * 1024 * 3]byte
		n, err := clientConn.Read(buf[:])
		if err != nil {
			log.Printf("[%s]Failed to read data from the client: %v", msgId, err)
			break
		}
		ws.Lock()
		err = ws.Conn.WriteMessage(websocket.BinaryMessage, (&common.Proto{
			MsgType: common.ReqData,
			Data:    buf[:n],
			MsgId:   msgId,
		}).ToBytes())
		ws.Unlock()
		if err != nil {
			log.Printf("[%s]Failed to write data to the WebSocket server: %v", msgId, err)
			break
		}
	}
	ws.Conn.Close()
	clientConn.Close()
	log.Printf("[%s]Closing the connection", msgId)
}
