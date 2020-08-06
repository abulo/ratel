#!/bin/bash

#拉取缺少的模块，移除不用的模块。
go mod tidy;
#下载依赖包
go mod download;
#打印模块依赖图
go mod graph;
#将依赖复制到vendor下
go mod vendor;
#校验依赖
go mod verify;
#解释为什么需要依赖
go mod why;