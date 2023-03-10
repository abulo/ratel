package lang

// Placeholder is a placeholder object that can be used globally.
var Placeholder PlaceholderType

// AnyType ...
type (
	// AnyType can be used to hold any type.
	AnyType = any
	// PlaceholderType represents a placeholder type.
	PlaceholderType = struct{}
)
