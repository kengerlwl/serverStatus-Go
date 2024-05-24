#!/bin/bash

# 创建输出目录
mkdir -p ./buildPackage

# 编译客户端为 Windows 可执行文件
echo "Compiling client for Windows..."
GOOS=windows GOARCH=amd64 go build -o ./buildPackage/client_probe.exe client.go config.go client_info.go
if [ $? -ne 0 ]; then
    echo "Failed to compile client for Windows."
    exit 1
fi

# 编译客户端为 Linux 可执行文件
echo "Compiling client for Linux..."
GOOS=linux GOARCH=amd64 go build -o ./buildPackage/client_linux_probe client.go config.go client_info.go
if [ $? -ne 0 ]; then
    echo "Failed to compile client for Linux."
    exit 1
fi

# 编译服务端为 Windows 可执行文件
echo "Compiling server for Windows..."
GOOS=windows GOARCH=amd64 go build -o ./buildPackage/server_probe.exe find.go config.go client_info.go
if [ $? -ne 0 ]; then
    echo "Failed to compile server for Windows."
    exit 1
fi

# 编译服务端为 Linux 可执行文件
echo "Compiling server for Linux..."
GOOS=linux GOARCH=amd64 go build -o ./buildPackage/server_linux_probe find.go config.go client_info.go
if [ $? -ne 0 ]; then
    echo "Failed to compile server for Linux."
    exit 1
fi

echo "Build completed successfully."
