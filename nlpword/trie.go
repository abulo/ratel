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
	Character rune
	// Children 为字符到子节点的映射；nil 表示尚无任何子节点（惰性分配，PERF-9-5）。
	Children map[rune]*Node
	// Failure 为 Aho-Corasick 失败指针。
	Failure *Node
	// Parent 为父节点，根节点的父节点为 nil。
	Parent *Node
	depth  int

	// failCache 缓存本节点在缺失字符上的压缩失败转移结果，避免查询期重复回溯（PERF-9-1，部分实施：惰性缓存等价于原 fail 遍历结果，重建失败链时清空）。
	failCache map[rune]*Node
	// outputLink 指向失败链上最近的结束节点（含自身），构建期计算，供 output 系列快速跳过非结束节点（PERF-9-2）。
	outputLink *Node
}

// BuildFailureLinks 按广度优先遍历重建整棵树的失败指针；须在词库构建完成后、匹配前调用。
func (tree *Trie) BuildFailureLinks() {
	queue := tree.queueBFS()
	for _, node := range queue {
		parent := node.Parent
		// PERF-9-4 [P2]：按标准 AC 规则利用父节点失败指针定位后继，替代从父节点逐层上溯；
		// 根节点直接子节点的失败指针指向根节点。
		if parent.IsRootNode() {
			node.Failure = parent
		} else {
			f := parent.Failure
			for f != nil && f != tree.Root && f.Children[node.Character] == nil {
				f = f.Failure
			}
			if f != nil {
				if sibling := f.Children[node.Character]; sibling != nil {
					node.Failure = sibling
				} else {
					node.Failure = tree.Root
				}
			} else {
				node.Failure = tree.Root
			}
		}

		// PERF-9-2 [P1] 预计算 output-link：当前节点本身结束则指向自身，
		// 否则继承失败链上最近的结束节点；等价于原 output 系列递归遍历结果。
		if node.IsPathEnd() {
			node.outputLink = node
		} else if f := node.Failure; f.IsPathEnd() {
			node.outputLink = f
		} else {
			node.outputLink = f.outputLink
		}

		// PERF-9-1 [P1] 失败链重建后，旧节点上缓存的压缩转移已失效，必须清空。
		node.failCache = nil
	}
}

// queueBFS 以类型化切片队列按广度优先顺序返回除根节点外的全部节点（PERF-9-4 [P2]：
// 替代 goroutine+channel 与 any 装箱链表队列，构建期不再产生 channel 收发与接口装箱分配）。
func (tree *Trie) queueBFS() []*Node {
	queue := make([]*Node, 0, len(tree.Root.Children))
	for _, child := range tree.Root.Children {
		queue = append(queue, child)
	}
	for i := 0; i < len(queue); i++ {
		n := queue[i]
		for _, child := range n.Children {
			queue = append(queue, child)
		}
	}
	return queue
}

// bfs 以通道方式按广度优先顺序逐个返回树中的节点，遍历结束后关闭通道。
//
// 保留通道签名以兼容包内既有测试；内部改用类型化切片队列实现（PERF-9-4 [P2]）。
func (tree *Trie) bfs() <-chan *Node {
	ch := make(chan *Node)
	go func() {
		defer close(ch)
		for _, n := range tree.queueBFS() {
			ch <- n
		}
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
			// PERF-9-5 [P2] 惰性分配子节点映射，避免叶子节点与尚未挂子节点的节点支付空 map 堆分配。
			if current.Children == nil {
				current.Children = make(map[rune]*Node)
			}
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
//
// PERF-9-3 [P1]（未实施）：报告指出此处失配回退存在 O(n²) 重扫，但该 trie 为
// 前缀过滤（未挂 AC 失败转移），改用失败指针回退可能改变被删除字符区间的边界，
// 无法在缺少等价样本的前提下证明结果一致；按语义兼容红线保留原回退实现，
// 待补充覆盖样本后单独评审。
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

	ac := new(ac)
	// PERF-9-2 [P1] 依据文本长度预设结果切片容量（封顶 256）；首次命中时才真正分配，
	// 保持原实现“无命中返回 nil 切片”的返回语义，减少命中密集场景的切片翻倍拷贝。
	if capHint := len(runes); capHint > 0 {
		if capHint > 256 {
			capHint = 256
		}
		ac.resultsCap = capHint
	}
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

// NewNode 创建一个代表指定字符的新节点；子节点映射保持 nil，首次插入子节点时再分配（PERF-9-5 [P2]）。
func NewNode(character rune) *Node {
	return &Node{
		Character: character,
	}
}

// NewRootNode 创建根节点，其失败指针指向自身。
func NewRootNode(character rune) *Node {
	root := &Node{
		isRootNode: true,
		Character:  character,
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
