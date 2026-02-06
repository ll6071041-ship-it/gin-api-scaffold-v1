#!/bin/sh

# 这是一个用于等待数据库启动的脚本
# 原始来源: https://github.com/eficode/wait-for/blob/master/wait-for

TIMEOUT=15
QUIET=0

echo "Waiting for $1..."

for i in $(seq $TIMEOUT) ; do
  # 使用 nc 命令检测端口 (你的 Dockerfile 已经安装了 netcat-openbsd)
  nc -z $(echo $1 | cut -d : -f 1) $(echo $1 | cut -d : -f 2) > /dev/null 2>&1
  
  result=$?
  if [ $result -eq 0 ] ; then
    if [ $QUIET -ne 1 ] ; then echo "$1 is available after $i seconds"; fi
    # 如果检测通了，就执行后面的命令 (-- 之后的内容)
    shift 2
    exec "$@"
    exit 0
  fi
  sleep 1
done

echo "Operation timed out" >&2
exit 1