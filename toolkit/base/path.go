package base

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"
)

func ratelHome() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	home := path.Join(dir, ".ratel")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func ratelHomeWithDir(dir string) string {
	home := path.Join(ratelHome(), dir)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

// textExtensions 为允许执行模块路径字面量替换的模板文本后缀；其余二进制/资源文件直接流式拷贝。
// PERF-8-11 [P2] 避免对图片、字体、压缩包等非文本文件整文件读入内存并逐替换全量扫描。
var textExtensions = map[string]struct{}{
	".go": {}, ".mod": {}, ".sum": {}, ".yaml": {}, ".yml": {}, ".json": {},
	".toml": {}, ".txt": {}, ".md": {}, ".proto": {}, ".sql": {}, ".vue": {},
	".ts": {}, ".js": {}, ".html": {}, ".css": {}, ".xml": {}, ".ini": {},
	".env": {}, ".sh": {},
}

// textBasenames 为无后缀但按文本处理的常见配置/清单文件名。
var textBasenames = map[string]struct{}{
	"Makefile": {}, "Dockerfile": {}, ".gitignore": {}, "README": {}, "LICENSE": {},
}

// replacePair 是一组预转换完成的 old/new 字节对；替换字面量在目录级只解析一次，
// 避免 copyFile 对每个文件重复 []byte(old)/[]byte(next) 分配（PERF-8-11 [P2]）。
type replacePair struct {
	oldB []byte
	newB []byte
}

// parseReplacePairs 将成对出现的替换字面量一次性转为字节切片对。
func parseReplacePairs(replaces []string) []replacePair {
	pairs := make([]replacePair, 0, len(replaces)/2)
	var old string
	for i, next := range replaces {
		if i%2 == 0 {
			old = next
			continue
		}
		pairs = append(pairs, replacePair{oldB: []byte(old), newB: []byte(next)})
	}
	return pairs
}

// isTextFile 判断目标文件是否属于需要做字面量替换的文本类型。
func isTextFile(name string) bool {
	if _, ok := textBasenames[filepath.Base(name)]; ok {
		return true
	}
	_, ok := textExtensions[strings.ToLower(filepath.Ext(name))]
	return ok
}

func copyFile(src, dst string, replaces []string) error {
	return copyFileWithPairs(src, dst, parseReplacePairs(replaces))
}

// copyFileWithPairs 使用目录级共享的替换字节对完成单文件拷贝；文本文件先判断字面量是否存在
// 再执行 ReplaceAll，非文本文件走流式 io.Copy，文件权限保持原文件模式。
func copyFileWithPairs(src, dst string, pairs []replacePair) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !isTextFile(src) {
		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcinfo.Mode())
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	}
	buf, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		// 未命中的替换直接跳过，避免无匹配时仍对整文件做一轮全量扫描与分配。
		if bytes.Contains(buf, pair.oldB) {
			buf = bytes.ReplaceAll(buf, pair.oldB, pair.newB)
		}
	}
	return os.WriteFile(dst, buf, srcinfo.Mode())
}

func copyDir(src, dst string, replaces, ignores []string) error {
	return copyDirWithPairs(src, dst, parseReplacePairs(replaces), ignores)
}

func copyDirWithPairs(src, dst string, pairs []replacePair, ignores []string) error {
	var err error
	var fds []os.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		if hasSets(fd.Name(), ignores) {
			continue
		}
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())
		var e error
		if fd.IsDir() {
			e = copyDirWithPairs(srcfp, dstfp, pairs, ignores)
		} else {
			e = copyFileWithPairs(srcfp, dstfp, pairs)
		}
		if e != nil {
			return e
		}
	}
	return nil
}

func hasSets(name string, sets []string) bool {
	return slices.Contains(sets, name)
}

// Tree 递归遍历目录并打印其中每个文件的相对路径与大小，用于在命令行中展示已生成的文件清单。
//
// dir 为打印相对路径时剥离的前缀目录；遍历过程中的错误被忽略，不中断输出。
func Tree(path string, dir string) {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			fmt.Printf("%s %s (%v bytes)\n", color.GreenString("CREATED"), strings.Replace(path, dir+"/", "", -1), info.Size())
		}
		return nil
	})
}
