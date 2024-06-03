#!/bin/bash

if [ -z "$1" ]; then
    binname=server
else
    binname=$1
fi

dir=`pwd`
echo $dir

#echo `ls -al`

cd ${dir}/server
# 设置临时环境变量
export GOPROXY=https://goproxy.cn,direct

go mod download
go mod tidy

# 打印go的版本
go version

# 编译程序
go build -o $binname

cp -afr $binname ../release
chmod u+x ../release/$binname


rm -rf $binname



