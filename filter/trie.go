package filter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"unicode/utf8"
)

// trie 表示自动机中的一个节点，保存子节点映射、失败指针、路径长度与结束标记。
type trie struct {
	// next 为字符到子节点的映射。
	next map[rune]*trie
	// fail 为失败指针。
	fail *trie
	// len 为该节点对应路径的字节长度。
	len  uint8
	// end 标记是否为单词结束节点。
	end  bool
}

// NewTrieWriter 创建一个空的词库写入器，内部初始化自动机根节点。
func NewTrieWriter() *TrieWriter {
	return &TrieWriter{tireRoot: &trie{next: map[rune]*trie{}}}
}

// TrieWriter 负责向 Trie 自动机插入词库、维护跳过字符集合并构建失败指针。
//
// 该类型不保证并发安全，词库构建应在单 goroutine 内完成后再交给扫描器使用。
type TrieWriter struct {
	// size 为词库中已插入的单词数量。
	size     int
	// skip 为匹配时需要跳过的字符集合。
	skip     *Skip
	// tireRoot 为自动机根节点。
	tireRoot *trie
}

// setSkip 设置需要跳过的字符集合并返回写入器本身。
func (t *TrieWriter) setSkip(skip *Skip) *TrieWriter {
	t.skip = skip
	return t
}

// Skip 返回当前跳过字符集合。
func (t *TrieWriter) Skip() *Skip {
	return t.skip
}

// Size 返回词库中已插入的单词数量。
func (t *TrieWriter) Size() int {
	return t.size
}

// Insert 向词库插入一个单词，自动跳过集合中的字符；返回写入器本身以支持链式调用。
func (t *TrieWriter) Insert(word string) *TrieWriter {
	node := t.tireRoot
	wLen := 0
	for _, v := range word {
		if t.skip.ShouldSkip(v) {
			continue
		}
		if _, ok := node.next[v]; !ok {
			node.next[v] = &trie{next: map[rune]*trie{}}
		}
		wLen += utf8.RuneLen(v)
		node = node.next[v]
		node.len = uint8(wLen)
	}
	if wLen > 0 && !node.end {
		node.end = true
		t.size++
	}
	return t
}

// InsertWords 批量插入词表中的全部单词，返回写入器本身以支持链式调用。
func (t *TrieWriter) InsertWords(words []string) *TrieWriter {
	for _, word := range words {
		t.Insert(word)
	}
	return t
}

// WriteString 将字符串作为单个单词插入词库，实现 io.Writer 接口；返回字符串长度与 nil 错误。
func (t *TrieWriter) WriteString(s string) (n int, err error) {
	t.Insert(s)
	return len(s), nil
}

// trie 返回自动机根节点。
func (t *TrieWriter) trie() *trie {
	return t.tireRoot
}

// InsertScanner 按扫描器每行读取单词并插入词库。
func (t *TrieWriter) InsertScanner(scanner *bufio.Scanner) {
	for scanner.Scan() {
		t.Insert(scanner.Text())
	}
}

// InsertFile 从指定文件按换行分隔读取单词并插入词库；文件打开失败时直接 panic。
func (t *TrieWriter) InsertFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()
	_, _ = t.InsertReader(file, '\n')
}

// InsertReader 从读取器按分隔符批量读取单词并插入词库，返回已读取的字节数与读取错误。
func (t *TrieWriter) InsertReader(reader io.Reader, delim byte) (n int, err error) {
	return t.insertReader(reader, make([]byte, 64*1024), delim)
}

func (t *TrieWriter) insertReader(reader io.Reader, buf []byte, delim byte) (n int, err error) {
	start, end, n1 := 0, 0, 0
	for {
		n1, err = reader.Read(buf[start:])
		if err == io.EOF {
			t.InsertBytes(buf[:end], delim)
			break
		}
		if err != nil {
			break
		}
		end = n1 + start
		last := end
		n += n1
		if last < len(buf) {
			t.InsertBytes(buf[:last], delim)
			break
		}
		last = bytes.LastIndexByte(buf[:end], delim)
		if last == -1 {
			err = errors.New("word too long")
			break
		}
		t.InsertBytes(buf[:last], delim)
		copy(buf, buf[last+1:end])
		start = end - last - 1
		end -= last + 1
	}
	return
}

// InsertBytes 按分隔符从字节切片解析单词并插入词库，插入过程中跳过集合中的字符；返回切片长度。
func (t *TrieWriter) InsertBytes(p []byte, delim byte) (n int) {
	n = len(p)
	for i := 0; i < n; {
		node := t.trie()
		wLen := 0

		for i < n && p[i] != delim {
			r, l := decodeBytes(p[i:])

			if t.skip.ShouldSkip(r) {
				i += l
				continue
			}

			if _, ok := node.next[r]; !ok {
				node.next[r] = &trie{next: map[rune]*trie{}}
			}
			wLen += l
			node = node.next[r]
			node.len = uint8(wLen)
			i += l
		}

		if wLen > 0 && !node.end {
			node.end = true
			t.size++
		}

		for i < n && p[i] == delim {
			i++
		}
	}

	return
}

// BuildFail 按广度优先遍历为所有节点构建失败指针，返回遍历的节点数量；
// 须在插入词库完成后、扫描前调用。
func (t *TrieWriter) BuildFail() int {
	queue := make([]*trie, 0, len(t.tireRoot.next))
	for _, node := range t.tireRoot.next {
		node.fail = t.tireRoot
		queue = append(queue, node)
	}

	count := 0
	for level := 1; len(queue) > 0; level++ {
		temp := make([]*trie, len(queue))
		copy(temp, queue)
		queue = queue[:0]
		for _, prev := range temp {
			for c, curr := range prev.next {
				count++
				queue = append(queue, curr)
				failTo := prev.fail
				for failTo != nil && failTo.next[c] == nil {
					failTo = failTo.fail
				}
				if failTo == nil {
					curr.fail = t.tireRoot
				} else {
					curr.fail = failTo.next[c]
				}
			}
		}
	}
	return count
}

// String 以每行一个单词的形式输出词库内容，最多输出 1000 条，避免词库过大时输出失控。
func (t *TrieWriter) String() string {
	type pair struct {
		trie  *trie
		runes []rune
	}
	buf := bytes.Buffer{}
	limit := 1000
	queue := []*pair{{t.trie(), nil}}
loop:
	for len(queue) > 0 {
		temp := make([]*pair, len(queue))
		copy(temp, queue)
		queue = queue[:0]
		for _, p := range temp {
			if p.trie.end {
				buf.WriteString(string(p.runes))
				buf.WriteByte('\n')
				limit--
				if limit == 0 {
					break loop
				}
			}
			for r, node := range p.trie.next {
				queue = append(queue, &pair{node, appendRunes(p.runes, r)})
			}
		}
	}
	return string(buf.Bytes()[:buf.Len()-1])
}

func appendRunes(runes []rune, r rune) []rune {
	res := make([]rune, len(runes)+1)
	copy(res, runes)
	res[len(res)-1] = r
	return res
}

// Array 返回词库中全部单词的字符串列表。
func (t *TrieWriter) Array() []string {
	type pair struct {
		trie  *trie
		runes []rune
	}
	res := make([]string, 0, t.Size())
	queue := []*pair{{t.trie(), nil}}
	for len(queue) > 0 {
		temp := make([]*pair, len(queue))
		copy(temp, queue)
		queue = queue[:0]
		for _, p := range temp {
			if p.trie.end {
				res = append(res, string(p.runes))
			}
			for r, node := range p.trie.next {
				queue = append(queue, &pair{node, appendRunes(p.runes, r)})
			}
		}
	}
	return res
}
