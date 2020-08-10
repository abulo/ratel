#!/bin/bash
# 程序所在目录
path=$(cd `dirname $0`; pwd)
for file in `ls $path`; do
    if [ -d $path"/"$file ];then
        if [ "vendor"!=$file ];then
            goimports -w $path"/"$file
            gofmt -w $path"/"$file
        fi
    fi
done