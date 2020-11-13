#!/bin/bash
Color_Text()
{
    echo -e " \e[0;$2m$1\e[0m"
}

Echo_Red()
{
    echo -e "\n"
    echo $(Color_Text "$1" "31")
    echo -e "\n"
}

Echo_Green()
{
    echo -e "\n"
    echo $(Color_Text "$1" "32")
    echo -e "\n"
}

Echo_Yellow()
{
    echo -e "\n"
    echo -n $(Color_Text "$1" "33")
    echo -e "\n"
}

Echo_Blue()
{
    echo -e "\n"
    echo $(Color_Text "$1" "34")
    echo -e "\n"
}

#https://colobu.com/2017/06/27/Lint-your-golang-code-like-a-mad-man/

# #gotype会对go文件和包进行语义(semantic)和句法(syntactic)的分析,这是google提供的一个工具。
# Echo_Red "==================================gotype分析========================================"
# Echo_Red "gotype会对go文件和包进行语义(semantic)和句法(syntactic)的分析"

# find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -not -path "./static/*"  -not -path  "./docker/*" -not -path  "./"  -not -path "./static" -print0 | xargs -0 gotype
# Echo_Red "==================================gotype完成========================================"


# #列出了所有复杂度大于12的函数
# Echo_Green "==================================gocyclo分析========================================"
# Echo_Green "列出了所有复杂度大的函数"
# gocyclo -over 12 $(ls -d */ | grep -v vendor)
# Echo_Green "==================================gocyclo完成========================================"

# #这个工具提供接口类型的建议，换句话说，它会对可以本没有必要定义成具体的类型的代码提出警告
# Echo_Yellow "==================================interfacer分析========================================"
# Echo_Yellow "这个工具提供接口类型的建议，换句话说，它会对可以本没有必要定义成具体的类型的代码提出警告"

# interfacer $(go list ./... | grep -v /vendor/)

# Echo_Yellow "==================================interfacer完成========================================"

# #deadcode会告诉你哪些代码片段根本没用
# Echo_Blue "==================================deadcode分析========================================"
# Echo_Blue "deadcode会告诉你哪些代码片段根本没用"
# find . -type d -not -path "./vendor/*" -not -path "./.git/*" -not -path "./static/*"  -not -path  "./docker/*" -not -path  "./"  -not -path "./static" | xargs deadcode
# Echo_Blue "==================================deadcode完成========================================"

#misspell用来拼写检查
Echo_Red "==================================misspell分析========================================"
Echo_Red "misspell用来拼写检查"
find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -not -path "./static/*"  -not -path  "./docker/*"   -not -path "./static"  -print0 | xargs -0 misspell
Echo_Red "==================================misspell完成========================================"

#goconst 会查找重复的字符串，这些字符串可以抽取成常量
# Echo_Green "==================================goconst分析========================================"
# Echo_Green "goconst 会查找重复的字符串，这些字符串可以抽取成常量"
# find . -type d -not -path "./vendor/*" -not -path "./.git/*" -not -path "./static/*"  -not -path  "./docker/*"  -not -path  "./" -not -path "./static"  | xargs goconst
# Echo_Green "==================================goconst完成========================================"

#静态检查
Echo_Yellow "==================================staticcheck分析========================================"
Echo_Yellow "静态检查"
staticcheck $(go list ./...);
Echo_Yellow "==================================staticcheck完成========================================"