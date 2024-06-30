# 第一阶段：构建Golang应用
FROM golang:1.22 AS builder

WORKDIR /app

# 拷贝应用代码到镜像中
COPY . .

# 编译应用
RUN go build -o cinetool

# 第二阶段：构建最终镜像
FROM debian:latest

# 安装SQLite
RUN apt-get update && apt-get install -y sqlite3

WORKDIR /app

# 从第一阶段复制编译好的应用到最终镜像
COPY --from=builder /app/cinetool .
COPY --from=builder /app/dist dist

# 设置容器启动时运行的命令
CMD ["./cinetool"]
