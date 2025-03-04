# 构建过程
FROM golang:1.22.3 AS builder

WORKDIR /app
COPY . .

# 设置Go环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

# 执行构建
RUN go env \
    && go mod tidy \
    && go build -o island ./cmd/main.go

# 执行过程
FROM alpine:latest

# 安装时区数据
RUN apk add --no-cache tzdata

# 设置工作目录
WORKDIR /app

# 复制构建好的二进制文件和配置文件
COPY --from=builder /app/island ./
COPY --from=builder /app/config.yaml ./

# 设置时区（例如，设置为上海时区）
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 暴露端口
EXPOSE 8080

# 设置入口点
ENTRYPOINT ["./island"]