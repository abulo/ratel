package nlpword

import (
	"strings"
	"sync"
)

type Node struct {
	val  rune
	path string

	depth     int
	term      bool //是否是node尾节点
	termCount int

	mask uint64
	meta interface{}

	parent   *Node
	children map[rune]*Node
}

type Trie struct {
	lock sync.Mutex

	root *Node
	size int
}

func NewTrie() *Trie {
	return &Trie{
		root: &Node{depth: 0, parent: nil, children: make(map[rune]*Node)},
		size: 0,
	}
}

func (t *Trie) Root() *Node {
	return t.root
}

/**
 * 前缀树添加数据
 *
 * @return 返回当前数据的尾节点
 */
func (t *Trie) AddKey(key string, meta interface{}) *Node {
	t.lock.Lock()
	t.size++

	runes := []rune(key)
	mask := maskval(runes)

	node := t.root
	node.mask |= mask
	node.termCount++
	for i := range runes {
		r := runes[i]
		mask = maskval(runes[i:])
		if n, ok := node.children[r]; ok {
			node = n
			node.mask |= mask
		} else {
			node = node.NewChild(r, "", mask, nil, false)
		}
		node.termCount++
	}
	node = node.NewChild(0x0, key, 0, meta, true) //为数据node创建尾节点，尾节点包含数据完整信息
	t.lock.Unlock()

	return node //返回的是尾节点
}

func (t *Trie) FindKey(key string) (*Node, bool) {
	node := findNode(t.Root(), []rune(key))
	if node == nil {
		return nil, false
	}

	//找到node的尾节点
	node, ok := node.Children()[0x0]
	if !ok || !node.term {
		return nil, false
	}

	return node, true
}

func (t *Trie) RemoveKey(key string) {
	var (
		i    int
		rs   = []rune(key)
		node = findNode(t.Root(), rs)
	)

	t.lock.Lock()

	t.size--
	for n := node.Parent(); n != nil; n = n.Parent() {
		i++
		if len(n.Children()) > 1 {
			r := rs[len(rs)-i]
			n.RemoveChild(r)
			break
		}
	}

	t.lock.Unlock()
}

/**
 * 通过数据的前缀搜索所有匹配的数据key
 *
 */
func (t *Trie) SearchKeyByPre(pre string) []string {
	node := findNode(t.Root(), []rune(pre))
	if node == nil {
		return nil
	}

	return nodeReach(node)
}

/**
 * 通过数据的前缀和后缀找到确定的一个key
 *
 */
func (t *Trie) SearchKeyByPre2Next(pre, next string) (string, bool) {
	if len(pre) == 0 && len(next) == 0 {
		return ``, false
	}
	if len(pre) == 0 && len(next) > 0 {
		return next, true
	}
	if len(pre) > 0 && len(next) == 0 {
		return pre, true
	}

	pnode := findNode(t.Root(), []rune(pre))
	if pnode == nil {
		return ``, false
	}

	if pnode.termCount < 1 {
		return ``, false
	}

	key, ok := nodeReachByPre2Next(pnode, pre, next)
	if ok {
		return key, ok
	}

	return ``, false
}

/**
 * 通过数据的前缀和一个后缀查找节点，如果存在节点则返回前缀+后缀的字符串
 *
 */
func (t *Trie) SearchKeyNodeByPre2Next(pre, next string) (string, bool) {
	if len(pre) == 0 && len(next) == 0 {
		return ``, false
	}
	if len(pre) == 0 && len(next) > 0 {
		return next, true
	}
	if len(pre) > 0 && len(next) == 0 {
		return pre, true
	}

	ok := findNodeByPre2Next(t.Root(), []rune(pre), []rune(next))
	if ok {
		return pre + next, ok
	}

	return ``, false
}

func (t *Trie) HasKey(key string) bool {
	node := findNode(t.Root(), []rune(key))
	return node != nil
}

func (n *Node) Parent() *Node {
	return n.parent
}

func (n *Node) Children() map[rune]*Node {
	return n.children
}

func (n *Node) Val() rune {
	return n.val
}

func (n *Node) Depth() int {
	return n.depth
}

func (n *Node) Mask() uint64 {
	return n.mask
}

func (parent *Node) NewChild(val rune, path string, mask uint64, meta interface{}, term bool) *Node {
	node := &Node{
		val:      val,
		path:     path,
		mask:     mask,
		meta:     meta,
		term:     term,
		parent:   parent,
		children: make(map[rune]*Node),
		depth:    parent.depth + 1,
	}

	parent.children[node.val] = node
	parent.mask |= mask
	return node
}

func (n *Node) RemoveChild(r rune) {
	delete(n.children, r)
	for np := n.parent; np != nil; np = np.parent {
		np.mask ^= np.mask
		np.mask |= uint64(1) << uint64(np.val-'a')
		for _, npc := range np.children {
			np.mask |= npc.mask
		}
	}
}

/**
 * 通过字符rune找到字符rune最后的node
 *
 */
func findNode(node *Node, runes []rune) *Node {
	if node == nil {
		return nil
	}

	if len(runes) == 0 {
		return node
	}

	n, ok := node.Children()[runes[0]]
	if !ok {
		return nil
	}

	var nrunes []rune
	if len(runes) > 1 {
		nrunes = runes[1:]
	} else {
		nrunes = runes[0:0]
	}

	return findNode(n, nrunes)
}

/**
 * 通过前缀字符查找后一个字符的节点，判断是否存在
 *
 */
func findNodeByPre2Next(node *Node, pre []rune, next []rune) bool {
	if node == nil {
		return false
	}

	if len(pre) == 0 || len(next) != 1 {
		return false
	}

	n, ok := node.Children()[pre[0]]
	if !ok {
		return false
	}

	var rpre []rune
	if len(pre) > 1 {
		rpre = pre[1:]
	} else if len(pre) == 1 {
		nextn, ok := n.Children()[next[0]]
		if !ok {
			return false
		}

		if nextn.Val() != next[0] {
			return false
		}

		return true
	}

	return findNodeByPre2Next(n, rpre, next)
}

func maskval(rd []rune) uint64 {
	var m uint64
	for _, d := range rd {
		m |= uint64(1) << uint64(d-'a')
	}

	return m
}

func nodeReach(node *Node) []string {
	var n *Node
	var i int

	keys := make([]string, 0)
	nodes := make([]*Node, 1)
	nodes[0] = node

	l := len(nodes)
	for ; l != 0; l = len(nodes) {
		i = l - 1
		n = nodes[i]
		nodes = nodes[:i]
		for _, child := range n.Children() {
			nodes = append(nodes, child)
		}
		if n.term {
			word := n.path
			keys = append(keys, word)
		}
	}

	return keys
}

func nodeReachByPre2Next(node *Node, pre, next string) (string, bool) {
	var n *Node
	var i int

	rlpre := len([]rune(pre))
	nodes := make([]*Node, 1)
	nodes[0] = node

	l := len(nodes)
	for ; l != 0; l = len(nodes) {
		i = l - 1
		n = nodes[i]
		nodes = nodes[:i]
		for _, child := range n.Children() {
			nodes = append(nodes, child)
		}
		if n.term {
			rword := []rune(n.path)
			nextkey := string(rword[rlpre:])
			if strings.Compare(next, nextkey) == 0 {
				return n.path, true
			}
		}
	}

	return ``, false
}
