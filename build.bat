set dir=%~dp0
cd %dir%

go version
set GOPROXY=https://goproxy.cn,direct
go build

