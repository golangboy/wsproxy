FROM golang as builder
WORKDIR /app
COPY . .
WORKDIR /app/client
RUN export CGO_ENABLED=0 && go build -o app .
WORKDIR /app/server
RUN export CGO_ENABLED=0 && go build -o app .
WORKDIR /app
RUN chmod +x docker-entrypoint.sh

FROM alpine
COPY --from=builder /app /app
EXPOSE 80 1180
CMD ["/app/docker-entrypoint.sh"]
