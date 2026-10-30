package nlpword

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
)

// Filter 是基于 Trie 自动机的敏感词过滤器。
//
// 零值不可直接使用，须通过 New 构造；该过滤器不保证并发安全，多个 goroutine
// 同时加载词库或执行检测替换时需要外部加锁。
type Filter struct {
	trie       *Trie
	noise      *regexp.Regexp
	buildVer   int64
	updatedVer int64
}

// New 创建一个空的敏感词过滤器，默认噪声正则匹配空白与常用干扰符号。
func New() *Filter {
	return &Filter{
		trie:  NewTrie(),
		noise: regexp.MustCompile(`[\|\s&%$@*]+`),
	}
}

// UpdateNoisePattern 更换噪声过滤正则表达式；传入非法正则时该方法会 panic。
func (filter *Filter) UpdateNoisePattern(pattern string) {
	filter.noise = regexp.MustCompile(pattern)
}

// LoadWordDict 从本地字典文件按行加载敏感词；文件打开或读取失败时返回错误。
func (filter *Filter) LoadWordDict(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			logger.Logger.Error("Error closing filter: ", err)
		}
	}()

	return filter.Load(f)
}

// LoadNetWordDict 从指定 URL 下载字典并按行加载敏感词；请求超时 5 秒，
// 网络请求或内容读取失败时返回错误。
func (filter *Filter) LoadNetWordDict(url string) error {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	rsp, err := c.Get(url)
	if err != nil {
		return err
	}

	defer func() {
		if err := rsp.Body.Close(); err != nil {
			logger.Logger.Error("Error closing rsp: ", err)
		}
	}()

	return filter.Load(rsp.Body)
}

// Load 从读取器按行读取敏感词并加入词库；读取过程出错时返回错误，EOF 为正常结束。
func (filter *Filter) Load(rd io.Reader) error {
	buf := bufio.NewReader(rd)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		filter.AddWord(string(line))
	}

	return nil
}

// updateFailureLink 在词库发生变更后按需重建自动机失败指针，避免每次匹配前重复构建。
func (filter *Filter) updateFailureLink() {
	if filter.buildVer != filter.updatedVer {
		filter.trie.BuildFailureLinks()
		filter.buildVer = filter.updatedVer
	}
}

// AddWord 向词库追加一个或多个敏感词；追加后下一次匹配前会自动重建失败指针。
func (filter *Filter) AddWord(words ...string) {
	filter.trie.Add(words...)
	filter.updatedVer = time.Now().UnixNano()
}

// Filter 删除文本中全部命中的敏感词，返回过滤后的文本。
func (filter *Filter) Filter(text string) string {
	filter.updateFailureLink()
	return filter.trie.Filter(text)
}

// Replace 将文本中全部命中的敏感词替换为 repl 字符，返回替换后的文本。
func (filter *Filter) Replace(text string, repl rune) string {
	filter.updateFailureLink()
	return filter.trie.Replace(text, repl)
}

// FindIn 检测文本是否包含敏感词；命中时返回 true 与第一个敏感词，否则返回 false 与空字符串。
//
// 检测前会先按噪声正则去除干扰字符。
func (filter *Filter) FindIn(text string) (bool, string) {
	filter.updateFailureLink()
	text = filter.RemoveNoise(text)
	return filter.trie.FindIn(text)
}

// FindAll 返回文本中命中的全部敏感词列表；未命中时返回空列表。
func (filter *Filter) FindAll(text string) []string {
	filter.updateFailureLink()
	return filter.trie.FindAll(text)
}

// Validate 校验文本是否合法；合法返回 true 与空字符串，不合法返回 false 与第一个命中的敏感词。
//
// 检测前会先按噪声正则去除干扰字符。
func (filter *Filter) Validate(text string) (bool, string) {
	filter.updateFailureLink()
	text = filter.RemoveNoise(text)
	return filter.trie.Validate(text)
}

// RemoveNoise 按当前噪声正则删除文本中的干扰字符后返回结果文本。
func (filter *Filter) RemoveNoise(text string) string {
	return filter.noise.ReplaceAllString(text, "")
}
