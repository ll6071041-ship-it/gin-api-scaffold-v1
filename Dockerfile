# 1. 编译阶段
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# 2. 运行阶段 (必须用 debian 安装 netcat)
FROM debian:bookworm-slim

WORKDIR /app

# 安装 netcat (nc)，这是 wait-for.sh 必须的工具
RUN apt-get update && apt-get install -y netcat-openbsd && rm -rf /var/lib/apt/lists/*

# 复制二进制文件
COPY --from=builder /app/main .
# 复制配置文件
COPY --from=builder /app/config.yaml .
# 复制 wait-for 脚本
COPY --from=builder /app/wait-for.sh .

# 给脚本可执行权限 (非常重要！)
RUN chmod +x ./wait-for.sh

# 声明端口
EXPOSE 8080

# 注意：这里不写 CMD，因为我们在 docker-compose.yaml 里写了 command 覆盖它