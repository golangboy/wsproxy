FROM golang
WORKDIR /app
COPY . .
WORKDIR /app/client
RUN go build -o app .
WORKDIR /app/server
RUN go build -o app .
WORKDIR /app
RUN chmod +x start.sh
EXPOSE 80 1180
CMD ["/app/start.sh"]
