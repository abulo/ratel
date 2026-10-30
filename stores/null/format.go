package null

// 日期时间格式常量，供本包各时间类型的序列化与解析使用。
const (
	// RFC3339DateOnly 为仅日期格式 YYYY-MM-DD。
	RFC3339DateOnly = "2006-01-02"
	// RFC3339TimeOnly 为仅时间格式 hh:mm:ss。
	RFC3339TimeOnly = "15:04:05"
	// DateTimeSQL 为 SQL 日期时间格式 YYYY-MM-DD hh:mm:ss。
	DateTimeSQL = "2006-01-02 15:04:05"
	// TimeStampSQL 为带微秒的 SQL 时间戳格式 YYYY-MM-DD hh:mm:ss.000000。
	TimeStampSQL = "2006-01-02 15:04:05.000000"
)
