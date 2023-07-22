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

# 🧑‍💻 How to use
## On server
```shell
./server -listen=":80"
```

## On client
```shell
./client -listen=":1180" -ws="your_server_ip:80"
```
### 💻 CLI Tool
```shell
export all_proxy=socks5://localhost:1180
curl https://google.com -v
```
### 🐱 Clash
```yaml
proxies:
  - name: "wsproxy"
    type: socks5
    server: localhost
    port: 1180
```
> add `wsproxy` to your proxy group