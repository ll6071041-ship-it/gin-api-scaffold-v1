# 定义变量
BINARY_NAME=bluebell
GO_FILES=$(shell find . -name "*.go")

# 伪目标：防止文件夹里有叫 build、run 的文件导致冲突
.PHONY: all build run test clean tool lint help build-linux build-windows

# 默认目标：直接输入 make 时执行的动作
all: build

# 1. 编译
build:
	@echo "正在编译 $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go

# 2. 运行
run:
	@echo "正在运行..."
	go run main.go

# 3. 测试
test:
	@echo "正在运行单元测试..."
	go test -v ./...

# 4. 清理 (删除二进制文件)
clean:
	@echo "正在清理..."
	@if [ -f $(BINARY_NAME) ] ; then rm $(BINARY_NAME) ; fi
	@if [ -f $(BINARY_NAME).exe ] ; then rm $(BINARY_NAME).exe ; fi

# 5. 格式化代码 & 整理依赖
tidy:
	go fmt ./...
	go mod tidy

# 6. 交叉编译 - 编译出 Linux 可用的二进制文件 (常用于部署)
build-linux:
	@echo "正在编译 Linux 版本..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux main.go

# 7. 交叉编译 - 编译出 Windows 可用的 exe (如果你在 Mac/Linux 上开发)
build-windows:
	@echo "正在编译 Windows 版本..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe main.go

# 8. 帮助信息
help:
	@echo "make build - 编译当前平台版本"
	@echo "make run   - 直接运行"
	@echo "make clean - 清理构建文件"
	@echo "make tidy  - 格式化代码并整理 mod"
	@echo "make build-linux - 编译 Linux 版本"