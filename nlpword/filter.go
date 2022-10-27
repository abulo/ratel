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

// Filter 敏感词过滤器
type Filter struct {
	trie       *Trie
	noise      *regexp.Regexp
	buildVer   int64
	updatedVer int64
}

// New 返回一个敏感词过滤器
func New() *Filter {
	return &Filter{
		trie:  NewTrie(),
		noise: regexp.MustCompile(`[\|\s&%$@*]+`),
	}
}

// UpdateNoisePattern 更新去噪模式
func (filter *Filter) UpdateNoisePattern(pattern string) {
	filter.noise = regexp.MustCompile(pattern)
}

// LoadWordDict 加载敏感词字典
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

// LoadNetWordDict 加载网络敏感词字典
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

// Load common method to add words
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

func (filter *Filter) updateFailureLink() {
	if filter.buildVer != filter.updatedVer {
		// fmt.Println("update failure link")
		filter.trie.BuildFailureLinks()
		filter.buildVer = filter.updatedVer
	}
}

// AddWord 添加敏感词
func (filter *Filter) AddWord(words ...string) {
	filter.trie.Add(words...)
	filter.updatedVer = time.Now().UnixNano()
}

// Filter 过滤敏感词
func (filter *Filter) Filter(text string) string {
	filter.updateFailureLink()
	return filter.trie.Filter(text)
}

// Replace 和谐敏感词
func (filter *Filter) Replace(text string, repl rune) string {
	filter.updateFailureLink()
	return filter.trie.Replace(text, repl)
}

// FindIn 检测敏感词
func (filter *Filter) FindIn(text string) (bool, string) {
	filter.updateFailureLink()
	text = filter.RemoveNoise(text)
	return filter.trie.FindIn(text)
}

// FindAll 找到所有匹配词
func (filter *Filter) FindAll(text string) []string {
	filter.updateFailureLink()
	return filter.trie.FindAll(text)
}

// Validate 检测字符串是否合法
func (filter *Filter) Validate(text string) (bool, string) {
	filter.updateFailureLink()
	text = filter.RemoveNoise(text)
	return filter.trie.Validate(text)
}

// RemoveNoise 去除空格等噪音
func (filter *Filter) RemoveNoise(text string) string {
	return filter.noise.ReplaceAllString(text, "")
}
