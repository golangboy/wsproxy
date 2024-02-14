# 📖 Introduction
A lightweight proxy tool based on the WebSocket protocol
# 🚀 Features
- Lightweight
- No configuration
- Wrap CDN
- Support SOCKS5/HTTP
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
# 🐳 Docker

## From Build
```shell
git clone github.com/golangboy/wsproxy
cd wsproxy
docker build -t wsproxy .
```
#### client
```shell
docker run -itd -p 1180:1180 -e ws=your_server:80 wsproxy
```
#### server
```shell
docker run -itd -p 80:80 wsproxy
```

## From DockerHub
```shell
docker pull golangboyme/wsproxy
```
#### client
```shell
docker run -itd -p 1180:1180 -e ws=your_server:80 golangboyme/wsproxy
```
#### server
```shell
docker run -itd -p 80:80 golangboyme/wsproxy
```
# 💻 Test
```shell
export all_proxy=socks5://localhost:1180
curl https://google.com -v
```