#!/bin/bash

# 自定义二进制消息结构生成
python3 go2ts.py -i ../comm/command.go -o go 

if [ "$1" != "svn" ]
then
    exit
fi

svn commit protos/ -m "$2"


