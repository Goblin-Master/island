# 构建阶段
FROM golang:1.22.3 AS builder

# 设置工作目录
WORKDIR /app

# 将代码复制到容器中
COPY . .

# 设置环境变量
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.baidu.com,direct

# 缓存依赖项（确保 go.mod 和 go.sum 在 COPY 之后）
COPY go.mod go.sum ./
RUN go mod download  # 下载依赖项到本地缓存

# 构建项目
RUN go build -o island ./cmd/main.go

# 执行阶段
FROM alpine:latest

# 安装时区数据
RUN apk add --no-cache tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件和配置文件
COPY --from=builder /app/island ./
COPY --from=builder /app/config.yaml ./

# 设置时区（例如，设置为上海时区）
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 暴露端口
EXPOSE 8080

# 设置入口点
ENTRYPOINT ["./island"]