package lang

// Placeholder 是可全局复用的空结构体占位实例，不携带任何数据。
var Placeholder PlaceholderType

type (
	// AnyType 表示可容纳任意类型的值，等价于 any。
	AnyType = any
	// PlaceholderType 表示不携带数据的空结构体占位类型。
	PlaceholderType = struct{}
)
