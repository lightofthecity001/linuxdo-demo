# Stage 1: 构建阶段 (Builder)
FROM golang:alpine AS builder

# 设置环境变量：开启 Go Modules，配置国内代理 (关键步骤)
ENV GO111MODULE=on \
    #GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置工作目录
WORKDIR /build

# 先只拷贝 go.mod 和 go.sum，利用 Docker 缓存层加速下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝所有源代码
COPY . .

# 编译项目，生成名为 app 的二进制文件
RUN go build -o app .

# -------------------------------------------

# Stage 2: 运行阶段 (Runner)
# 使用极小的 Alpine Linux 镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# (可选) 安装时区数据，防止日志时间差8小时
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

# 从构建阶段把编译好的二进制文件“偷”过来
COPY --from=builder /build/app .

# 暴露端口 (和你的 Gin 监听端口一致)
EXPOSE 10110

# 启动命令
CMD ["./app"]