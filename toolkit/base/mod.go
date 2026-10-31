package base

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/mod/modfile"
)

// moduleVersionCache 缓存 ModuleVersion 按模块路径查询成功的结果；
// PERF-8-13 [P2] 同一命令进程内重复查询时不再重复执行 go mod graph 全图读取与逐行扫描。
// 失败结果不写入缓存，保留原有失败后再次尝试的语义。（预期收益需实测验证）
var moduleVersionCache sync.Map

// ModulePath 读取指定 go.mod 文件并返回其 module 声明路径；文件读取失败时返回错误。
func ModulePath(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}

// ModuleVersion 通过 go mod graph 在模块依赖图中查找指定模块的路径与版本；未找到或命令失败时返回错误。
func ModuleVersion(path string) (string, error) {
	if cached, ok := moduleVersionCache.Load(path); ok {
		return cached.(string), nil
	}
	stdout := &bytes.Buffer{}
	fd := exec.Command("go", "mod", "graph")
	fd.Stdout = stdout
	fd.Stderr = stdout
	if err := fd.Run(); err != nil {
		return "", err
	}
	rd := bufio.NewReader(stdout)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			return "", err
		}
		str := string(line)
		i := strings.Index(str, "@")
		if strings.Contains(str, path+"@") && i != -1 {
			result := path + str[i:]
			moduleVersionCache.Store(path, result)
			return result, nil
		}
	}
}

// RatelMod 返回当前环境中 ratel 框架源码或模块缓存的本地目录路径；无法获取版本时回退到 GOPATH 源码目录。
func RatelMod() string {
	cacheOut, _ := exec.Command("go", "env", "GOMODCACHE").Output()
	cachePath := strings.Trim(string(cacheOut), "\n")
	pathOut, _ := exec.Command("go", "env", "GOPATH").Output()
	gopath := strings.Trim(string(pathOut), "\n")
	if cachePath == "" {
		cachePath = filepath.Join(gopath, "pkg", "mod")
	}
	if path, err := ModuleVersion("github.com/abulo/ratel/v3"); err == nil {
		return filepath.Join(cachePath, path)
	}
	return filepath.Join(gopath, "src", "github.com", "abulo", "ratel")
}
