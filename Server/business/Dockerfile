# 使用官方Go运行时作为基础镜像
FROM golang:latest AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用alpine作为最终镜像
FROM alpine:latest

# 安装ca-certificates
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从builder镜像复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/business-config.json .

# 暴露端口
EXPOSE 8081

# 运行应用
CMD ["./main"]
