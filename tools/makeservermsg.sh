#!/bin/bash

# 自定义二进制消息结构生成
python3 go2go.py -i ../servercomm/command.go -o go 
go fmt ../servercomm/command_binary.go
