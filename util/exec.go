package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"
)

// Exec 解析 command 并执行，命令支持用双引号或单引号包裹包含空格的参数。
//
// 执行成功时将按行拆分的全部输出写入 output，returnVar 置 0，并返回最后一行输出；
// 执行失败时 returnVar 置 1 并返回空字符串。
func Exec(command string, output *[]string, returnVar *int) string {
	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		*returnVar = 1
		return ""
	}
	*returnVar = 0
	*output = strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	if l := len(*output); l > 0 {
		return (*output)[l-1]
	}
	return ""
}

// PERF-7-31 [P2] 带缓存上限的输出写入器：未触顶时写入 stdBuf，触顶后仅返回写入成功、不再缓存；终端透传由 MultiWriter 另一路承担
type boundedCacheWriter struct {
	buf    *bytes.Buffer
	limit  int64
	capped atomic.Bool
}

func (w *boundedCacheWriter) Write(p []byte) (int, error) {
	if !w.capped.Load() {
		if int64(w.buf.Len())+int64(len(p)) > w.limit {
			remain := w.limit - int64(w.buf.Len())
			if remain > 0 {
				w.buf.Write(p[:remain])
			}
			w.capped.Store(true)
		} else {
			w.buf.Write(p)
		}
	}
	return len(p), nil
}

// System 解析 command 并执行，标准输出与标准错误同时透传到终端并缓存到内存。
//
// 执行成功时 returnVar 置 0，返回缓存输出的最后一行；启动、读取或等待失败时
// returnVar 置 1，错误信息打印到标准输出并返回空字符串。
func System(command string, returnVar *int) string {
	*returnVar = 0
	var stdBuf bytes.Buffer
	var err, err1, err2, err3 error

	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	// PERF-7-31 [P2] stdBuf 设置 1MB 缓存上限，超出后只向终端透传、不再缓存，避免无界缓冲随长命令输出无限增长
	cache := &boundedCacheWriter{buf: &stdBuf, limit: 1 << 20}
	stdout := io.MultiWriter(os.Stdout, cache)
	stderr := io.MultiWriter(os.Stderr, cache)

	err = cmd.Start()
	if err != nil {
		*returnVar = 1
		return ""
	}

	// PERF-7-31 [P2] 两个拷贝 goroutine 纳入 WaitGroup，cmd.Wait() 后先等待拷贝结束再读取错误，避免等待返回后缓冲丢失与悬挂 goroutine
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		defer wg.Done()
		_, err2 = io.Copy(stderr, stderrIn)
	}()

	err3 = cmd.Wait()
	wg.Wait()
	if err1 != nil || err2 != nil || err3 != nil {
		if err1 != nil {
			fmt.Println(err1)
		}
		if err2 != nil {
			fmt.Println(err2)
		}
		if err3 != nil {
			fmt.Println(err3)
		}
		*returnVar = 1
		return ""
	}
	if output := strings.TrimRight(stdBuf.String(), "\n"); output != "" {
		pos := strings.LastIndex(output, "\n")
		if pos == -1 {
			return output
		}
		return output[pos+1:]
	}
	return ""
}

// PassThru 解析 command 并执行，标准输出与标准错误直接透传到终端，不做缓存。
//
// 执行成功时 returnVar 置 0；失败时 returnVar 置 1，并将错误打印到标准输出。
func PassThru(command string, returnVar *int) {
	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		*returnVar = 1
		fmt.Println(err)
	} else {
		*returnVar = 0
	}
}
