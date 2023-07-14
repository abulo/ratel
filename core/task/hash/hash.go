package hash

import "hash/crc32"

// Hash 返回数据的哈希值
func Hash(data []byte) uint32 { return crc32.ChecksumIEEE(data) }
