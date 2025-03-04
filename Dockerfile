# 构建阶段
FROM golang:1.22.3 AS builder

# 设置工作目录
WORKDIR /app

# 设置环境变量
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.tencent.com,direct
ENV GODEBUG=x509ignoreCN

# 安装 CA 证书
RUN apk add --no-cache ca-certificates

# 复制代码到容器中
COPY . .

# 缓存依赖项
COPY go.mod go.sum ./
RUN go mod download

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

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 暴露端口
EXPOSE 8080

# 设置入口点
ENTRYPOINT ["./island"]