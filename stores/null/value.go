package null

// Value 描述可空类型的状态查询接口，用于区分未赋值、赋值为 null 与持有有效值三种状态。
type Value interface {
	// IsValid 报告是否持有显式赋值且非空的取值。
	IsValid() bool
	// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
	IsSet() bool
}
