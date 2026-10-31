package nlpword

// ac 是 Aho-Corasick 匹配过程中的临时状态，负责沿失败指针回溯并收集匹配结果。
type ac struct {
	results    []string
	resultsCap int
}

// fail 沿失败指针回溯，直到找到字符 c 的子节点；到达根节点仍无匹配时返回根节点。
//
// PERF-9-1 [P1]（部分实施，预期收益需实测验证）：对 (节点, 字符) 的失败转移结果做惰性缓存，
// 重复失配字符命中缓存后 O(1) 返回；缓存目标与原逐层回溯结果完全一致，失败链重建时由
// BuildFailureLinks 清空。完整 goto 转移表属于算法级改造，未在本批次实施，建议单独评审。
func (ac *ac) fail(node *Node, c rune) *Node {
	if cached, ok := node.failCache[c]; ok {
		return cached
	}

	target := node
	cur := node.Failure
	for {
		if child, ok := cur.Children[c]; ok {
			target = child
			break
		}
		if cur.IsRootNode() {
			target = cur
			break
		}
		cur = cur.Failure
	}

	if node.failCache == nil {
		node.failCache = make(map[rune]*Node)
	}
	node.failCache[c] = target
	return target
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
//
// PERF-9-2 [P1]：基于构建期预计算的 output-link 迭代遍历失败链上的结束节点，
// 不再递归走完整条失败链；收集顺序与原递归实现一致。首次命中时按预设容量惰性分配，
// 保持原实现“无命中返回 nil 切片”的返回语义。
func (ac *ac) output(node *Node, runes []rune, position int) {
	for cur := node.outputLink; cur != nil; cur = cur.Failure.outputLink {
		if ac.results == nil {
			ac.results = make([]string, 0, ac.resultsCap)
		}
		ac.results = append(ac.results, string(runes[position+1-cur.depth:position+1]))
	}
}

// firstOutput 返回沿失败指针回溯命中的第一个敏感词，未命中时返回空字符串。
//
// PERF-9-2 [P1]：直接通过 output-link 定位最近的结束节点，跳过失败链上的非结束节点；
// 返回的命中片段与原递归实现一致。
func (ac *ac) firstOutput(node *Node, runes []rune, position int) string {
	if cur := node.outputLink; cur != nil {
		return string(runes[position+1-cur.depth : position+1])
	}
	return ""
}

// replace 将当前位置命中的敏感词在 rune 切片中改写为指定字符，并沿失败指针继续处理重叠命中。
//
// PERF-9-2 [P1]：递归改为等价迭代循环，避免长失败链上的栈递归调用；改写范围与原实现一致。
func (ac *ac) replace(node *Node, runes []rune, position int, replace rune) {
	for cur := node; !cur.IsRootNode(); cur = cur.Failure {
		if cur.IsPathEnd() {
			for i := position + 1 - cur.depth; i < position+1; i++ {
				runes[i] = replace
			}
		}
	}
}
