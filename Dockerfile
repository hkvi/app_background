# 多阶段构建Dockerfile - 使用已有镜像
# 第一阶段：构建阶段
FROM golang:latest AS builder

# 设置工作目录
WORKDIR /app

# 设置Go模块代理（使用国内镜像源）
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 复制go mod文件
COPY code/go.mod code/go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY code/ ./

# 构建应用
RUN go build -a -installsuffix cgo -ldflags="-w -s" -o main .

# 第二阶段：运行阶段 - 使用scratch最小镜像
FROM scratch

# 复制CA证书（用于HTTPS请求）
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/config.json .

# 暴露端口
EXPOSE 8080

# 健康检查（使用内置的HTTP客户端）
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/main", "-health-check"] || exit 1

# 启动应用
CMD ["./main"]
