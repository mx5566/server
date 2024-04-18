@echo off
set dir=%~dp0
cd %dir%/server
@rem 设置环境变量
set GOPROXY=https://goproxy.cn,direct

@rem 打印go的版本
go version

@rem 编译程序
go build

copy  /y server.exe  ..\release

del /F /Q  server.exe



