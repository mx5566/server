#!/bin/bash
dir=`pwd`
echo $dir

cd $dir/server
# 设置临时环境变量
export GOPROXY=https://goproxy.cn,direct

# 打印go的版本
go version

# 编译程序
go build


cp -afr server ../release

rm -rf server



