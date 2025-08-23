@echo off
REM 远程部署脚本

REM 1. 打包后端
cd backend
set GOOS=linux
set GOARCH=amd64
REM go build -o card-authorization-app
cd ..


REM 2. 上传到服务器
REM 请替换下面的用户名和服务器IP
set SERVER_USER=root
set SERVER_IP=wangxiang-pro.top

REM 上传后端
scp ./backend/card-authorization-app %SERVER_USER%@%SERVER_IP%:/home/card/

REM 上传前端
scp -r frontend %SERVER_USER%@%SERVER_IP%:/home/card/

REM 3. 远程部署
REM 进入服务器执行部署脚本 /home/cardAuthorizationApp.sh
ssh %SERVER_USER%@%SERVER_IP% "bash /home/card/cardAuthorizationApp.sh"

echo "远程部署成功"