package nlpword

// LinkList 是基于链表实现的先进先出队列，供 Trie 树广度优先遍历使用。
//
// 该队列不保证并发安全，仅在单次遍历的同一 goroutine 内使用。
type LinkList struct {
	head  *listNode
	tail  *listNode
	count int64
}

// Push 在队尾追加一个元素。
func (list *LinkList) Push(v any) {
	node := &listNode{
		Value: v,
	}
	if list.head == nil {
		list.head = node
	} else {
		list.tail.Next = node

	}
	list.tail = node
	list.count++
}

// Pop 取出队首元素并返回其值；队列为空时返回 nil。
func (list *LinkList) Pop() any {
	if list.Empty() {
		return nil
	}

	n := list.head
	list.head = n.Next
	list.count--
	return n.Value
}

// Empty 判断队列是否为空。
func (list *LinkList) Empty() bool {
	return list.count == 0
}

type listNode struct {
	Value any
	Next  *listNode
}
