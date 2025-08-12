@echo off
cd /d "%~dp0backend"
echo 正在启动功能卡片授权系统...
echo 请确保已安装Go 1.21或更高版本
echo.
echo 正在下载依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo 依赖下载失败，请检查网络连接
    pause
    exit /b 1
)
echo.
echo 正在启动服务器...
go run main.go
pause
