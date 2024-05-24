@echo off
setlocal

REM 创建输出目录
if not exist .\buildPackage (
    mkdir .\buildPackage
)

REM 编译客户端为 Windows 可执行文件
echo Compiling client for Windows...
go build -o .\buildPackage\client_probe.exe client.go config.go client_info.go
if %errorlevel% neq 0 (
    echo Failed to compile client for Windows.
    exit /b 1
)

REM 设置环境变量以编译客户端为 Linux 可执行文件
echo Compiling client for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o .\buildPackage\client_linux_probe client.go config.go client_info.go
if %errorlevel% neq 0 (
    echo Failed to compile client for Linux.
    exit /b 1
)

REM 编译服务端为 Windows 可执行文件
echo Compiling server for Windows...
set GOOS=windows
set GOARCH=amd64
go build -o .\buildPackage\server_probe.exe find.go config.go client_info.go
if %errorlevel% neq 0 (
    echo Failed to compile server for Windows.
    exit /b 1
)

REM 设置环境变量以编译服务端为 Linux 可执行文件
echo Compiling server for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o .\buildPackage\server_linux_probe find.go config.go client_info.go
if %errorlevel% neq 0 (
    echo Failed to compile server for Linux.
    exit /b 1
)

echo Build completed successfully.

endlocal
