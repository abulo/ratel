// Package null 提供一组可空基础类型，用于替代 sql.NullXxx 处理数据库与 JSON 中的空值。
//
// 每个类型由内部取值、Valid（是否持有非空有效值）与 Set（是否被显式赋值，包含赋值为 null）
// 三个状态位组成，并统一实现 sql.Scanner、driver.Valuer、json.Marshaler、json.Unmarshaler、
// encoding.TextMarshaler 与 encoding.TextUnmarshaler 接口；空值在 JSON 中序列化为 null，
// 读写数据库时映射为 SQL NULL。
package null
