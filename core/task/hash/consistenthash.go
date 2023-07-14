package hash

import (
	"sort"
	"strconv"
)

const minReplicas = 3

type (
	// Func defines the hash method.
	Func func(data []byte) uint32

	// ConsistentHash is a ring hash implementation.
	ConsistentHash struct {
		hashFunc Func
		replicas int
		hashMap  map[uint32]string
		keys     []uint32 // Sorted
	}
)

// NewConsistentHash returns a ConsistentHash.
func NewConsistentHash() *ConsistentHash {
	return NewCustomConsistentHash(minReplicas, Hash)
}

// NewCustomConsistentHash returns a ConsistentHash with given replicas and hash func.
func NewCustomConsistentHash(replicas int, fn Func) *ConsistentHash {
	if replicas < minReplicas {
		replicas = minReplicas
	}

	if fn == nil {
		fn = Hash
	}

	return &ConsistentHash{
		hashFunc: fn,
		replicas: replicas,
		hashMap:  make(map[uint32]string),
	}
}

// IsEmpty Returns true if there are no items available.
func (h *ConsistentHash) IsEmpty() bool {
	return len(h.keys) == 0
}

func (h *ConsistentHash) Add(keys ...string) *ConsistentHash {
	for _, key := range keys {
		for i := 0; i < h.replicas; i++ {
			hash := h.hashFunc([]byte(key + strconv.Itoa(i)))

			if h.hashMap[hash] == "" {
				h.keys = append(h.keys, hash)
				h.hashMap[hash] = key
			}
		}
	}

	// sort
	sort.Slice(h.keys, func(i, j int) bool { return h.keys[i] < h.keys[j] })

	return h
}

// Get the closest item in the hash to the provided key.
func (h *ConsistentHash) Get(key string) string {
	if h.IsEmpty() {
		return ""
	}

	hash := h.hashFunc([]byte(key))

	index := sort.Search(len(h.keys), func(i int) bool { return h.keys[i] >= hash }) % len(h.keys)

	return h.hashMap[h.keys[index]]
}
