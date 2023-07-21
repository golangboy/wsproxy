# 📖 Introduction
A lightweight proxy tool based on the WebSocket protocol
# 🚀 Features
- Lightweight
- No configuration
- Wrap CDN
# 🔨️ Build
```shell
git clone github.com/golangboy/wsproxy
```

## Client
```shell
cd client
go build .
./client -h
```

```shell
Usage of ./client:
  -listen string
        listen socks5 address (default ":1180")
  -ws string
        websocket server address (default "localhost:80")

```

### Server
```shell
cd server
go build .
./server -h
```

```shell
Usage of ./server:
  -listen string
        listen address (default ":80")

```