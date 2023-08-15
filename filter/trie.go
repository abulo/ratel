package filter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"unicode/utf8"
)

// trie表示trie树中的每个节点，具有next、fail、len和end四个属性。
type trie struct {
	next map[rune]*trie // 映射表，用于存储下一个字符的节点
	fail *trie          // 指向该节点的失败指针
	len  uint8          // 表示该节点代表的字符串长度
	end  bool           // 表示是否是一个单词的结尾节点
}

// NewTrieWriter 返回一个新的TrieWriter对象，其中tireRoot属性为一个空的trie树根节点。
func NewTrieWriter() *TrieWriter {
	return &TrieWriter{tireRoot: &trie{next: map[rune]*trie{}}}
}

// TrieWriter 表示一个Trie写入器，包含了一些trie树的相关操作。
type TrieWriter struct {
	size     int   // trie树中单词的数量
	skip     *Skip // 需要跳过的字符集合
	tireRoot *trie // trie树根节点
}

// setSkip设置需要跳过的字符集合，并返回当前对象。
func (t *TrieWriter) setSkip(skip *Skip) *TrieWriter {
	t.skip = skip
	return t
}

// Skip 获取需要跳过的字符集合。
func (t *TrieWriter) Skip() *Skip {
	return t.skip
}

// Size 获取trie树中单词的数量。
func (t *TrieWriter) Size() int {
	return t.size
}

// Insert 向trie树中插入一个单词，返回当前对象。
func (t *TrieWriter) Insert(word string) *TrieWriter {
	node := t.tireRoot // 从trie树的根节点开始
	wLen := 0
	for _, v := range word { // 遍历单词中的每个字符
		if t.skip.ShouldSkip(v) { // 如果该字符应该跳过，则继续下一个字符
			continue
		}
		if _, ok := node.next[v]; !ok { // 如果下一个节点不存在，则创建一个新节点
			node.next[v] = &trie{next: map[rune]*trie{}}
		}
		wLen += utf8.RuneLen(v) // 更新单词长度
		node = node.next[v]
		node.len = uint8(wLen)
	}
	if wLen > 0 && !node.end { // 如果单词不为空，将结尾节点标记为end，并且数量加1。
		node.end = true
		t.size++
	}
	return t
}

// InsertWords 向trie树中插入一个字符串数组中的所有单词，返回当前对象。调用了Insert(word string)方法。
func (t *TrieWriter) InsertWords(words []string) *TrieWriter {
	for _, word := range words {
		t.Insert(word)
	}
	return t
}

// WriteString 将一个字符串写入到trie树中，返回写入的字节数和nil错误。调用了Insert(word string)方法。
func (t *TrieWriter) WriteString(s string) (n int, err error) {
	t.Insert(s)
	return len(s), nil
}

// trie返回trie树的根节点。
func (t *TrieWriter) trie() *trie {
	return t.tireRoot
}

func (t *TrieWriter) InsertScanner(scanner *bufio.Scanner) {
	for scanner.Scan() {
		t.Insert(scanner.Text())
	}
}

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

// InsertBytes 将一个字节数组写入到trie树中，返回写入的字节数和nil错误。在遍历字节数组的过程中，跳过被定义在skip属性中的字符，如果遇到换行符则在该单词的结尾节点标记为end
func (t *TrieWriter) InsertBytes(p []byte, delim byte) (n int) {
	n = len(p)           // 获取字节数组的长度
	for i := 0; i < n; { // 遍历字节数组
		node := t.trie() // 从trie树的根节点开始
		wLen := 0

		for i < n && p[i] != delim { // 判断字符是否为换行符并且还没有遍历完整个字节数组
			r, l := decodeBytes(p[i:]) // 解码字节数组中的一个rune，并且获取该rune的字节数量

			// 跳过一些无意义的字符
			if t.skip.ShouldSkip(r) {
				i += l
				continue
			}

			if _, ok := node.next[r]; !ok { // 如果下一个节点不存在，则创建一个新节点
				node.next[r] = &trie{next: map[rune]*trie{}}
			}
			wLen += l // 更新单词长度
			node = node.next[r]
			node.len = uint8(wLen) // 更新节点代表的字符串长度
			i += l                 // 向后移动光标
		}

		if wLen > 0 && !node.end { // 如果单词不为空，将结尾节点标记为end，并且数量加1。
			node.end = true
			t.size++
		}

		for i < n && p[i] == delim { // 判断是否为分隔符
			i++ // 向后移动光标
		}
	}

	return
}

// BuildFail 用于构建trie树中每个节点的失败指针，返回当前对象。
func (t *TrieWriter) BuildFail() int {
	// 创建队列并将根节点的所有子节点加入队列中
	queue := make([]*trie, 0, len(t.tireRoot.next))
	for _, node := range t.tireRoot.next {
		node.fail = t.tireRoot // 设置节点的失败指针为根节点
		queue = append(queue, node)
	}

	count := 0                                // 计数器用于记录遍历的节点数
	for level := 1; len(queue) > 0; level++ { // 遍历节点直到队列为空
		temp := make([]*trie, len(queue)) // 创建临时队列并将队列中的元素复制到临时队列中
		copy(temp, queue)
		queue = queue[:0]           // 清空队列
		for _, prev := range temp { // 遍历临时队列中的节点
			for c, curr := range prev.next { // 遍历当前节点的所有子节点
				count++
				queue = append(queue, curr)                  // 将子节点加入队列
				failTo := prev.fail                          // 获取当前节点的失败指针
				for failTo != nil && failTo.next[c] == nil { // 查找失败指针直到找到一个匹配字符c的节点或者到达根节点
					failTo = failTo.fail
				}
				if failTo == nil { // 如果没有找到匹配字符c的节点，则将当前节点的失败指针设置为根节点
					curr.fail = t.tireRoot
				} else { // 如果找到了匹配字符c的节点，则将当前节点的失败指针设置为该节点
					curr.fail = failTo.next[c]
				}
			}
		}
	}
	return count // 返回遍历的节点数
}

func (t *TrieWriter) String() string {
	type pair struct { // 定义一个键值对，trie表示trie树中的节点，runes表示从根节点到当前节点的字符序列
		trie  *trie
		runes []rune
	}
	buf := bytes.Buffer{}             // 创建缓冲区
	limit := 1000                     // 设置限制输出结果的条数
	queue := []*pair{{t.trie(), nil}} // 创建队列并将根节点加入队列中
loop:
	for len(queue) > 0 { // 遍历队列直到队列为空或者达到限制的输出结果数
		temp := make([]*pair, len(queue)) // 创建临时队列并将队列中的元素复制到临时队列中
		copy(temp, queue)
		queue = queue[:0]        // 清空队列
		for _, p := range temp { // 遍历临时队列中的节点
			if p.trie.end { // 如果当前节点是单词的结尾，则将字符序列转为字符串并添加到缓冲区中
				buf.WriteString(string(p.runes))
				buf.WriteByte('\n')
				limit--
				if limit == 0 { // 如果达到限制的输出结果数，则跳出循环
					break loop
				}
			}
			for r, node := range p.trie.next { // 遍历当前节点的所有子节点
				queue = append(queue, &pair{node, appendRunes(p.runes, r)}) // 将子节点加入队列，并更新字符序列
			}
		}
	}
	return string(buf.Bytes()[:buf.Len()-1]) // 返回缓冲区中的字符串，去掉最后一个换行符
}

func appendRunes(runes []rune, r rune) []rune {
	res := make([]rune, len(runes)+1)
	copy(res, runes)
	res[len(res)-1] = r
	return res
}

func (t *TrieWriter) Array() []string {
	type pair struct { // 定义一个键值对，trie表示trie树中的节点，runes表示从根节点到当前节点的字符序列
		trie  *trie
		runes []rune
	}
	res := make([]string, 0, t.Size())
	queue := []*pair{{t.trie(), nil}}
	for len(queue) > 0 { // 遍历队列直到队列为空或者达到限制的输出结果数
		temp := make([]*pair, len(queue)) // 创建临时队列并将队列中的元素复制到临时队列中
		copy(temp, queue)
		queue = queue[:0]        // 清空队列
		for _, p := range temp { // 遍历临时队列中的节点
			if p.trie.end { // 如果当前节点是单词的结尾，则将字符序列转为字符串并添加到缓冲区中
				res = append(res, string(p.runes))
			}
			for r, node := range p.trie.next { // 遍历当前节点的所有子节点
				queue = append(queue, &pair{node, appendRunes(p.runes, r)}) // 将子节点加入队列，并更新字符序列
			}
		}
	}
	return res
}
