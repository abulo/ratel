package nlpword

// Trie 是由敏感词构成的 Aho-Corasick 自动机树。
//
// 零值不可直接使用，须通过 NewTrie 构造；向树中插入新词后、匹配前须调用
// BuildFailureLinks 重建失败指针。该树不保证并发安全。
type Trie struct {
	// Root 为自动机根节点。
	Root *Node
}

// Node 表示 Trie 树上的一个节点，保存字符、子节点、失败指针、父节点与层级信息。
type Node struct {
	isRootNode bool
	isPathEnd  bool
	// Character 为该节点对应的字符。
	Character  rune
	// Children 为字符到子节点的映射。
	Children   map[rune]*Node
	// Failure 为 Aho-Corasick 失败指针。
	Failure    *Node
	// Parent 为父节点，根节点的父节点为 nil。
	Parent     *Node
	depth      int
}

// BuildFailureLinks 按广度优先遍历重建整棵树的失败指针；须在词库构建完成后、匹配前调用。
func (tree *Trie) BuildFailureLinks() {
	for node := range tree.bfs() {
		pointer := node.Parent
		var link *Node
		for link == nil {
			if pointer.IsRootNode() {
				link = pointer
				break
			}
			link = pointer.Failure.Children[node.Character]
			pointer = pointer.Failure

		}
		node.Failure = link

	}
}

// bfs 以通道方式按广度优先顺序逐个返回树中的节点，遍历结束后关闭通道。
func (tree *Trie) bfs() <-chan *Node {
	ch := make(chan *Node)
	go func() {
		queue := new(LinkList)
		for _, child := range tree.Root.Children {
			queue.Push(child)
		}

		for !queue.Empty() {
			n := queue.Pop().(*Node)
			ch <- n
			for _, child := range n.Children {
				queue.Push(child)
			}
		}

		close(ch)
	}()
	return ch
}

// NewTrie 创建一棵空 Trie 树并初始化根节点。
func NewTrie() *Trie {
	return &Trie{
		Root: NewRootNode(0),
	}
}

// Add 向树中插入一个或多个敏感词。
func (tree *Trie) Add(words ...string) {
	for _, word := range words {
		tree.add(word)
	}
}

func (tree *Trie) add(word string) {
	var current = tree.Root
	var runes = []rune(word)
	for position := range runes {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := NewNode(r)
			newNode.depth = current.depth + 1
			newNode.Parent = current
			current.Children[r] = newNode
			current = newNode
		}
		if position == len(runes)-1 {
			current.isPathEnd = true
		}
	}
}

// Replace 将文本中全部命中的敏感词替换为 character 字符，返回替换后的文本。
func (tree *Trie) Replace(text string, character rune) string {
	var (
		node  = tree.Root
		next  *Node
		runes = []rune(text)
	)

	var ac = new(ac)
	for position := range runes {
		next = ac.next(node, runes[position])
		if next == nil {
			next = ac.fail(node, runes[position])
		}

		node = next
		ac.replace(node, runes, position, character)
	}

	return string(runes)
}

// Filter 删除文本中全部命中的敏感词，返回过滤后的文本。
func (tree *Trie) Filter(text string) string {
	var (
		parent      = tree.Root
		current     *Node
		left        = 0
		found       bool
		runes       = []rune(text)
		length      = len(runes)
		resultRunes = make([]rune, 0, length)
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found {
			resultRunes = append(resultRunes, runes[left])
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() {
			left = position + 1
		}
		parent = current
	}

	resultRunes = append(resultRunes, runes[left:]...)
	return string(resultRunes)
}

// Validate 校验文本是否合法；不合法时返回 false 与第一个命中的敏感词，
// 合法时返回 true 与空字符串。
func (tree *Trie) Validate(text string) (bool, string) {
	const EMPTY = ""
	var (
		node  = tree.Root
		next  *Node
		runes = []rune(text)
	)

	var ac = new(ac)
	for position := range runes {
		next = ac.next(node, runes[position])
		if next == nil {
			next = ac.fail(node, runes[position])
		}

		node = next
		if first := ac.firstOutput(node, runes, position); len(first) > 0 {
			return false, first
		}
	}

	return true, EMPTY
}

// FindIn 判断文本中是否包含词库中的词，返回是否命中与第一个命中的敏感词。
func (tree *Trie) FindIn(text string) (bool, string) {
	validated, first := tree.Validate(text)
	return !validated, first
}

// FindAll 返回文本中命中的全部敏感词列表；未命中时返回空列表。
func (tree *Trie) FindAll(text string) []string {
	var (
		node  = tree.Root
		next  *Node
		runes = []rune(text)
	)

	var ac = new(ac)
	for position := range runes {
		next = ac.next(node, runes[position])
		if next == nil {
			next = ac.fail(node, runes[position])
		}

		node = next
		ac.output(node, runes, position)
	}

	return ac.results

}

// NewNode 创建一个代表指定字符的新节点。
func NewNode(character rune) *Node {
	return &Node{
		Character: character,
		Children:  make(map[rune]*Node, 0),
	}
}

// NewRootNode 创建根节点，其失败指针指向自身。
func NewRootNode(character rune) *Node {
	root := &Node{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*Node, 0),
		depth:      0,
	}

	root.Failure = root

	return root
}

// IsLeafNode 判断节点是否为叶子节点（无子节点）。
func (node *Node) IsLeafNode() bool {
	return len(node.Children) == 0
}

// IsRootNode 判断节点是否为根节点。
func (node *Node) IsRootNode() bool {
	return node.isRootNode
}

// IsPathEnd 判断节点是否为某个敏感词路径的结束节点。
func (node *Node) IsPathEnd() bool {
	return node.isPathEnd
}
