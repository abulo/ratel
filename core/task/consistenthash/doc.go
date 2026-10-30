// Package consistenthash 提供一致性哈希环实现，通过虚拟副本在节点集合间均匀分散键，
// 并在节点增减时尽量缩小需要迁移的键范围。
package consistenthash
