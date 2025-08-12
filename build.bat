@echo off
cd /d "%~dp0backend"
echo 正在编译功能卡片授权系统...
echo 正在下载依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo 依赖下载失败，请检查网络连接
    pause
    exit /b 1
)
echo.
echo 正在编译...
go build -o ../card-authorization.exe main.go
if %errorlevel% neq 0 (
    echo 编译失败，请检查代码
    pause
    exit /b 1
)
echo.
echo 编译成功！
echo 可执行文件已生成：card-authorization.exe
echo 双击运行即可启动服务器
pause
