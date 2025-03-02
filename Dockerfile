# 构建过程
FROM golang:1.22.3 AS  builder

WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o island ./cmd/main.go


# 执行过程
FROM alpine:latest



WORKDIR /app

COPY --from=builder /app/island  ./
COPY --from=builder /app/config.yaml  ./

EXPOSE 8080
ENTRYPOINT ./island