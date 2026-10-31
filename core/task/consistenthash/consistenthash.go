package consistenthash

/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 根据输入字节串计算 32 位哈希值；New 未指定时默认使用 crc32.ChecksumIEEE。
type Hash func(data []byte) uint32

// Map 是一致性哈希环，虚拟副本与实际节点一一映射，用于在节点集合变化时尽量减少键的迁移。
//
// 零值不可直接使用，须通过 New 构造；Map 内部不加锁，Add 期间不得并发调用 Get。
type Map struct {
	hash     Hash
	replicas int
	// keys 为排序后的虚拟节点哈希列表。
	keys    []int
	hashMap map[int]string
}

// New 创建一致性哈希环；replicas 为每个节点的虚拟副本数，fn 为哈希函数，
// fn 为 nil 时使用 CRC32 IEEE 校验。
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// IsEmpty 报告哈希环中是否尚未添加任何节点。
func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

// Add 将给定节点加入哈希环，每个节点按副本数生成对应的虚拟节点。
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 副本序号与节点名拼接后参与哈希，避免节点名自带数字前缀造成分布偏移。
			hash := int(m.hash([]byte(strconv.Itoa(i) + "_" + key)))
			if m.hashMap[hash] == "" {
				m.keys = append(m.keys, hash)
				m.hashMap[hash] = key
			}
		}
	}
	sort.Ints(m.keys)
}

// Get 返回哈希环上最接近 key 哈希值的节点；环为空时返回空字符串，
// 目标哈希超出末尾虚拟节点时回绕到环首。
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}
	// PERF-3-13 [P2]：复核结论——导出的 Hash 接口签名为 func([]byte) uint32，标准库 crc32 摘要不提供
	// WriteString 入口，消除 string→[]byte 拷贝需要变更导出 API（违反兼容红线）或自行实现字符串版 CRC32；
	// 该分配为一次小对象、触发频率等于任务调度频率，收益低于改造风险，本条目不改代码。
	hash := int(m.hash([]byte(key)))
	// 二分查找定位不小于目标哈希的虚拟节点。
	idx := sort.Search(len(m.keys), func(i int) bool { return m.keys[i] >= hash })
	// 超出末尾时回绕到环首虚拟节点。
	if idx == len(m.keys) {
		idx = 0
	}
	return m.hashMap[m.keys[idx]]
}
