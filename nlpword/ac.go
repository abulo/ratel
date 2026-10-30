package nlpword

// ac 是 Aho-Corasick 匹配过程中的临时状态，负责沿失败指针回溯并收集匹配结果。
type ac struct {
	results []string
}

// fail 沿失败指针回溯，直到找到字符 c 的子节点；到达根节点仍无匹配时返回根节点。
func (ac *ac) fail(node *Node, c rune) *Node {
	var next *Node
	for {
		next = ac.next(node.Failure, c)
		if next == nil {
			if node.IsRootNode() {
				return node
			}
			node = node.Failure
			continue
		}

		return next
	}

}

// next 返回节点在字符 c 上的子节点，不存在时返回 nil。
func (ac *ac) next(node *Node, c rune) *Node {
	next, ok := node.Children[c]
	if ok {
		return next
	}
	return nil
}

// output 沿失败指针收集当前位置命中的全部敏感词到 results。
func (ac *ac) output(node *Node, runes []rune, position int) {
	if node.IsRootNode() {
		return
	}

	if node.IsPathEnd() {
		ac.results = append(ac.results, string(runes[position+1-node.depth:position+1]))
	}

	ac.output(node.Failure, runes, position)
}

// firstOutput 返回沿失败指针回溯命中的第一个敏感词，未命中时返回空字符串。
func (ac *ac) firstOutput(node *Node, runes []rune, position int) string {
	if node.IsRootNode() {
		return ""
	}

	if node.IsPathEnd() {
		return string(runes[position+1-node.depth : position+1])
	}

	return ac.firstOutput(node.Failure, runes, position)
}

// replace 将当前位置命中的敏感词在 rune 切片中改写为指定字符，并沿失败指针继续处理重叠命中。
func (ac *ac) replace(node *Node, runes []rune, position int, replace rune) {
	if node.IsRootNode() {
		return
	}

	if node.IsPathEnd() {
		for i := position + 1 - node.depth; i < position+1; i++ {
			runes[i] = replace
		}
	}
	ac.replace(node.Failure, runes, position, replace)
}
