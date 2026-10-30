package resource

import "bytes"

type (
	// BatchError 聚合批量操作过程中产生的多个错误，用于结束时统一返回；
	// 零值可直接使用，方法为并发不安全，调用方须自行保证同一实例不被并发写入。
	BatchError struct {
		errs errorArray
	}

	errorArray []error
)

// Add 追加错误到批量错误中，nil 错误会被忽略。
func (be *BatchError) Add(errs ...error) {
	for _, err := range errs {
		if err != nil {
			be.errs = append(be.errs, err)
		}
	}
}

// Err 返回代表已聚合全部错误的 error：无错误时返回 nil，仅有一个错误时返回该错误，
// 多个错误时返回包含全部错误的 errorArray。
func (be *BatchError) Err() error {
	switch len(be.errs) {
	case 0:
		return nil
	case 1:
		return be.errs[0]
	default:
		return be.errs
	}
}

// NotNil 返回是否已经聚合了至少一个错误。
func (be *BatchError) NotNil() bool {
	return len(be.errs) > 0
}

// Error 将数组内错误以换行符拼接为一行一个的多行文本。
func (ea errorArray) Error() string {
	var buf bytes.Buffer

	for i := range ea {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(ea[i].Error())
	}

	return buf.String()
}
