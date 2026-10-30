package util

import (
	"fmt"
	"math"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/spf13/cast"
)

// datePatterns 保存 PHP 日期格式字符到 Go time 布局字符串的两两映射，
// 供 StrToTime 与 Date 在解析和格式化前替换格式串使用。
var datePatterns = []string{
	"Y", "2006",
	"y", "06",

	"m", "01",
	"n", "1",
	"M", "Jan",
	"F", "January",

	"d", "02",
	"j", "2",

	"D", "Mon",
	"l", "Monday",

	"g", "3",
	"G", "15",
	"h", "03",
	"H", "15",

	"a", "pm",
	"A", "PM",

	"i", "04",
	"s", "05",

	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	"r", time.RFC1123Z,
}

// DateTimeParse 将 "HH:MM:SS" 格式的时间字符串解析为自当日零点起经过的总秒数；
// 无法读出三个整数字段时返回零值与错误。
func DateTimeParse(st string) (int, error) {
	var h, m, s int
	n, err := fmt.Sscanf(st, "%d:%d:%d", &h, &m, &s)
	if err != nil || n != 3 {
		return 0, err
	}
	return h*3600 + m*60 + s, nil
}

// FormatDuring 计算 t 到当前时刻的经过时长，格式化为 "天小时分钟秒.十分之一秒" 的
// 中文字符串；t 在未来时天数部分带负号，其余时分秒仍取正值。
func FormatDuring(t time.Time) string {
	const (
		Decisecond = 100 * time.Millisecond
		Day        = 24 * time.Hour
	)
	ts := time.Since(t)
	sign := time.Duration(1)
	if ts < 0 {
		sign = -1
		ts = -ts
	}
	ts += +Decisecond / 2
	d := sign * (ts / Day)
	ts = ts % Day
	h := ts / time.Hour
	ts = ts % time.Hour
	m := ts / time.Minute
	ts = ts % time.Minute
	s := ts / time.Second
	ts = ts % time.Second
	f := ts / Decisecond
	return fmt.Sprintf("%d天%d小时%d分钟%d.%d秒", d, h, m, s, f)
}

// GetHourDiffer 返回 endTime 相对 startTime 的 Unix 秒差。
//
// 仅当 startTime 早于 endTime 时返回两者差值；其他情况返回 0。
func GetHourDiffer(startTime, endTime any) int64 {
	var s int64
	t1 := cast.ToTime(startTime)
	t2 := cast.ToTime(endTime)
	if t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		s = diff
		return s
	} else {
		return s
	}
}

// TimeZone 返回进程当前使用的本地时区。
func TimeZone() *time.Location {
	return time.Local
}

// Time 返回当前时刻在本地时区下的 Unix 时间戳（秒）。
func Time() int64 {
	cstZone := TimeZone()
	return time.Now().In(cstZone).Unix()
}

// ToTime 返回给定 time.Time 的 Unix 时间戳（秒）。
func ToTime(now time.Time) int64 {
	return now.Unix()
}

// Duration 将时长字符串（如 "1h30m"）解析为 time.Duration；格式非法时 panic。
func Duration(str string) time.Duration {
	dur, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	return dur
}

// Now 返回本地时区下的当前时刻。
func Now() time.Time {
	cstZone := TimeZone()
	return time.Now().In(cstZone)
}

// Timestamp 返回当前时刻的毫秒级 Unix 时间戳。
func Timestamp() int64 {
	cstZone := TimeZone()
	return time.Now().In(cstZone).UnixNano() / int64(time.Millisecond)
}

// StrToTime 按 PHP 格式 format 解析 strtime，返回本地时区下的 Unix 时间戳（秒）。
//
// 格式字符支持 datePatterns 中定义的 PHP 日期字符；解析失败时返回 0 与错误。
func StrToTime(format, strtime string) (int64, error) {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	cstZone := TimeZone()
	t, err := time.ParseInLocation(format, strtime, cstZone)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// Date 按 PHP 格式 format 格式化时间；不传可变参数时使用当前时刻，时区取 TimeZone。
func Date(format string, ts ...time.Time) string {

	cstZone := TimeZone()
	replacer := strings.NewReplacer(datePatterns...)
	formats := replacer.Replace(format)
	t := time.Now()
	if len(ts) > 0 {
		t = ts[0]
	}

	res := t.In(cstZone).Format(formats)

	// PHP 的 G 表示 24 小时制且不带前导零，布局串 "15" 在整点前仍可能补出 "0"，
	// 这里手工剥掉结果首位的前导零。
	if format == "G" && res[:1] == "0" {
		return res[1:]
	}
	return res
}

// CheckDate 校验年、月、日是否构成合法的公历日期，并按闰年规则校验 2 月天数。
func CheckDate(month, day, year int) bool {
	if month < 1 || month > 12 || day < 1 || day > 31 || year < 1 || year > 32767 {
		return false
	}
	switch month {
	case 4, 6, 9, 11:
		if day > 30 {
			return false
		}
	case 2:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			if day > 29 {
				return false
			}
		} else if day > 28 {
			return false
		}
	}

	return true
}

// Sleep 阻塞当前协程 t 秒。
func Sleep(t int64) {
	time.Sleep(time.Duration(t) * time.Second)
}

// USleep 阻塞当前协程 t 微秒。
func USleep(t int64) {
	time.Sleep(time.Duration(t) * time.Microsecond)
}

// UnixTimeFormatDate 将可转换为时间的值格式化为 "YYYY-MM-DD HH:mm:ss"。
func UnixTimeFormatDate(str any) string {
	return Date("Y-m-d H:i:s", cast.ToTimeInDefaultLocation(cast.ToInt64(str), TimeZone()))
}

// FormatDateTime 将可转换为时间的值格式化为 "YYYY-MM-DD HH:mm:ss"；值为空时返回空字符串。
func FormatDateTime(str any) string {
	if Empty(str) {
		return ""
	}
	return Date("Y-m-d H:i:s", cast.ToTimeInDefaultLocation(str, TimeZone()))
}

// FormatDate 将可转换为时间的值格式化为 "YYYY-MM-DD"。
func FormatDate(str any) string {
	return Date("Y-m-d", cast.ToTimeInDefaultLocation(str, TimeZone()))
}

// ordinalInYear 返回给定日期在当年中的序数，1 月 1 日为 1。
func ordinalInYear(year int, month time.Month, day int) (dayNo int) {
	return DateToJulian(year, month, day) - DateToJulian(year, 1, 1) + 1
}

// ISOWeekday 返回给定日期的 ISO 8601 星期编号，周一为 1，周日为 7。
//
// 该编号与标准库 time.Weekday（周日为 0）的取值不同。
func ISOWeekday(year int, month time.Month, day int) (weekday int) {
	return DateToJulian(year, month, day)%7 + 1
}

// DateToJulian 将公历日期转换为儒略日数。
func DateToJulian(year int, month time.Month, day int) (dayNo int) {
	// 采用公历至儒略日数的标准换算公式。
	if month < 3 {
		year = year - 1
		month = month + 12
	}
	year = year + 4800

	return day + (153*(int(month)-3)+2)/5 + 365*year +
		year/4 - year/100 + year/400 - 32045
}

// JulianToDate 将儒略日数转换为公历年、月、日。
func JulianToDate(dayNo int) (year int, month time.Month, day int) {
	e := 4*(dayNo+1401+(4*dayNo+274277)/146097*3/4-38) + 3
	h := e%1461/4*5 + 2

	day = h%153/5 + 1
	month = time.Month((h/153+2)%12 + 1)
	year = e/1461 - 4716 + (14-int(month))/12

	return year, month, day
}

// startOffset 返回给定 ISO 周的周一相对于年初的天数偏移，偏移值可能为负。
func startOffset(y, week int) (offset int) {
	// 由 Tomohiko Sakamoto 星期算法推导而来，等价于 week*7 - ISOWeekday(y,1,4) - 3。
	y = y - 1
	return week*7 - (y+y/4-y/100+y/400+3)%7 - 4
}

// StartDate 返回给定 ISO 8601 年、周数对应的周一日期。
func StartDate(wyear, week int) (year int, month time.Month, day int) {
	return JulianToDate(
		DateToJulian(wyear, 1, 1) + startOffset(wyear, week))
}

// FromYearWeek 返回给定日期所属的 ISO 8601 年份与周数。
func FromYearWeek(year, month, day int) (wyear, week int) {
	return fromYearWeek(year, time.Month(month), day)
}

func fromYearWeek(year int, month time.Month, day int) (wyear, week int) {
	week = (ordinalInYear(year, month, day) - ISOWeekday(year, month, day) + 10) / 7
	if week < 1 {
		return fromYearWeek(year-1, 12, 31)
	}

	if week == 53 &&
		DateToJulian(StartDate(year+1, 1)) <= DateToJulian(year, month, day) {
		return year + 1, 1
	}
	return year, week
}

// NextMonthDate 返回日期 d 所在月份下一个月的第一天与最后一天，格式均为 Y-m-d。
func NextMonthDate(d string) (string, string) {
	now := cast.ToTimeInDefaultLocation(d, TimeZone())
	var nextMonth int
	var nextYear int
	year, month, _ := now.Date()

	if int(month) == 12 {
		nextMonth = 1
		nextYear = year + 1
	} else {
		nextMonth = int(month) + 1
		nextYear = year
	}

	currentLocation := now.Location()

	firstDate := time.Date(nextYear, time.Month(nextMonth), 1, 0, 0, 0, 0, currentLocation)
	lastDate := firstDate.AddDate(0, 1, -1)

	return Date("Y-m-d", firstDate), Date("Y-m-d", lastDate)
}

// PrevMonthDate 返回日期 d 所在月份上一个月的第一天与最后一天，格式均为 Y-m-d。
func PrevMonthDate(d string) (string, string) {
	now := cast.ToTimeInDefaultLocation(d, TimeZone())
	var prevMonth int
	var prevYear int
	year, month, _ := now.Date()

	if int(month) == 1 {
		prevMonth = 12
		prevYear = year - 1
	} else {
		prevMonth = int(month) - 1
		prevYear = year
	}

	currentLocation := now.Location()

	firstDate := time.Date(prevYear, time.Month(prevMonth), 1, 0, 0, 0, 0, currentLocation)
	lastDate := firstDate.AddDate(0, 1, -1)

	return Date("Y-m-d", firstDate), Date("Y-m-d", lastDate)
}

// MonthDate 返回日期 d 所在月份的第一天与最后一天，格式均为 Y-m-d。
func MonthDate(d string) (string, string) {
	now := cast.ToTimeInDefaultLocation(d, TimeZone())
	year, month, _ := now.Date()
	currentLocation := now.Location()
	firstDate := time.Date(year, month, 1, 0, 0, 0, 0, currentLocation)
	lastDate := firstDate.AddDate(0, 1, -1)
	return Date("Y-m-d", firstDate), Date("Y-m-d", lastDate)
}

// Month 返回以 d 对应月份的上一个月为起点的连续 12 个月，分别给出 Y-m 字符串列表
// 与 Ym 数字列表，两个列表下标一一对应。
func Month(d string) ([]string, []int64) {
	now := cast.ToTimeInDefaultLocation(d, TimeZone())
	prevTime := now.AddDate(-1, 1, 0)

	res := make([]string, 0)
	resInt64 := make([]int64, 0)
	res = append(res, Date("Y-m", prevTime))
	resInt64 = append(resInt64, cast.ToInt64(Date("Ym", prevTime)))

	for i := 1; i <= 11; i++ {
		prevNewTime := prevTime.AddDate(0, i, 0)
		res = append(res, Date("Y-m", prevNewTime))
		resInt64 = append(resInt64, cast.ToInt64(Date("Ym", prevNewTime)))
	}
	return res, resInt64
}

// Day 返回起止日期 s 与 e 之间的全部日期，分别给出 Y-m-d 字符串列表与 Ymd 数字列表，
// 两个列表下标一一对应。
func Day(s string, e string) ([]string, []int64) {
	start := cast.ToTimeInDefaultLocation(s, TimeZone())
	end := cast.ToTimeInDefaultLocation(e, TimeZone())
	total := int(math.Abs(float64(start.Sub(end).Hours() / 24)))
	res := make([]string, 0)
	resInt64 := make([]int64, 0)
	res = append(res, Date("Y-m-d", start))
	resInt64 = append(resInt64, cast.ToInt64(Date("Ymd", start)))
	compareEnd := cast.ToInt64(Date("Ymd", end))

	for i := 1; i <= total; i++ {
		prevNewDay := start.AddDate(0, 0, i)
		res = append(res, Date("Y-m-d", prevNewDay))
		compareTemp := cast.ToInt64(Date("Ymd", prevNewDay))
		resInt64 = append(resInt64, compareTemp)
		if compareTemp == compareEnd {
			break
		}
	}
	return res, resInt64
}

// Hour 返回起止时刻 s 与 e 之间的全部小时，分别给出 "Y-m-d H" 字符串列表与
// YmdH 数字列表，两个列表下标一一对应。
func Hour(s string, e string) ([]string, []int64) {
	start := cast.ToTimeInDefaultLocation(s, TimeZone())
	end := cast.ToTimeInDefaultLocation(e, TimeZone())
	total := int(math.Abs(float64(start.Sub(end).Hours()))) + 23
	res := make([]string, 0)
	resInt64 := make([]int64, 0)

	for i := 0; i <= total; i++ {
		t, _ := time.ParseDuration(cast.ToString(i) + "h")
		tempTime := start.Add(t)
		res = append(res, Date("Y-m-d H", tempTime))
		compareTemp := cast.ToInt64(Date("YmdH", tempTime))
		resInt64 = append(resInt64, compareTemp)
	}
	return res, resInt64
}

// ToWeekDay 将可转换为时间的值转换为中文星期名，取值为周日至周六。
func ToWeekDay(t any) string {
	weekday := [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	now := cast.ToTimeInDefaultLocation(t, TimeZone())
	var y, m, c, year, month, day uint16
	year, month, day = cast.ToUint16(Date("Y", now)), cast.ToUint16(Date("m", now)), cast.ToUint16(Date("d", now))
	if month >= 3 {
		m = month
		y = year % 100
		c = year / 100
	} else {
		m = month + 12
		y = (year - 1) % 100
		c = (year - 1) / 100
	}

	// 使用有符号整数运算，避免无符号类型在取模时出现负数问题。
	week := int(y) + int(y)/4 + int(c)/4 - 2*int(c) + (26*(int(m)+1))/10 + int(day) - 1
	which_week := (week%7 + 7) % 7
	return weekday[which_week]
}
