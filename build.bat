@echo off
set dir=%~dp0
cd %dir%
set GOPROXY=https://goproxy.cn,direct

go version

go build

