package util

import (
	"fmt"
	"math"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/spf13/cast"
)

// DateFormat pattern rules.
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", // A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

// timeZone 默认时区
// var timeZone *time.Location

// DateTimeParse 时间解析
func DateTimeParse(st string) (int, error) {
	var h, m, s int
	n, err := fmt.Sscanf(st, "%d:%d:%d", &h, &m, &s)
	if err != nil || n != 3 {
		return 0, err
	}
	return h*3600 + m*60 + s, nil
}

// FormatDuring 格式化秒
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

// GetHourDiffer 获取相差时间
func GetHourDiffer(startTime, endTime interface{}) int64 {
	var s int64
	t1 := cast.ToTime(startTime)
	t2 := cast.ToTime(endTime)
	if t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		s = diff
		return s
	} else {
		return s
	}
}

// TimeZone GetTimeZone 获取时区
func TimeZone() *time.Location {
	// if timeZone == nil {
	// 	timeZone, _ = time.LoadLocation("Local")
	// }
	return time.Local
}

// // SetTimeZone 设置时区
// func SetTimeZone(zone string) *time.Location {
// 	//loc, err := time.LoadLocation("Asia/Shanghai")
// res, _ := time.LoadLocation(zone)
// 	timeZone = res

// 	fmt.Println(timeZone.String())
// 	return res
// }

// DateTime Functions

// Time time()
func Time() int64 {
	//cstZone := time.FixedZone("CST", 8*3600)
	cstZone := TimeZone()
	return time.Now().In(cstZone).Unix()
}

// Duration ...
func Duration(str string) time.Duration {
	dur, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	return dur
}

// Now time.Now()
func Now() time.Time {
	cstZone := TimeZone()
	return time.Now().In(cstZone)
}

// Timestamp 毫秒
func Timestamp() int64 {
	cstZone := TimeZone()
	return time.Now().In(cstZone).UnixNano() / int64(time.Millisecond)
}

// StrToTime StrToTime()
// StrToTime("02/01/2006 15:04:05", "02/01/2016 15:04:05") == 1451747045
// StrToTime("3 04 PM", "8 41 PM") == -62167144740
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

// Date date()
func Date(format string, ts ...time.Time) string {

	cstZone := TimeZone()
	replacer := strings.NewReplacer(datePatterns...)
	formats := replacer.Replace(format)
	t := time.Now()
	if len(ts) > 0 {
		t = ts[0]
	}

	res := t.In(cstZone).Format(formats)

	if format == "G" && res[:1] == "0" {
		return res[1:]
	}
	return res
}

// CheckDate Checkdate checkdate()
// Validate a Gregorian date
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
		// leap year
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

// Sleep sleep()
func Sleep(t int64) {
	time.Sleep(time.Duration(t) * time.Second)
}

// USleep Usleep usleep()
func USleep(t int64) {
	time.Sleep(time.Duration(t) * time.Microsecond)
}

// UnixTimeFormatDate ...
func UnixTimeFormatDate(str interface{}) string {
	return Date("Y-m-d H:i:s", cast.ToTimeInDefaultLocation(cast.ToInt64(str), TimeZone()))
}

// FormatDateTime ...
func FormatDateTime(str interface{}) string {
	if Empty(str) {
		return ""
	}
	return Date("Y-m-d H:i:s", cast.ToTimeInDefaultLocation(str, TimeZone()))
}

// FormatDate ...
func FormatDate(str interface{}) string {
	return Date("Y-m-d", cast.ToTimeInDefaultLocation(str, TimeZone()))
}

// ordinalInYear returns the ordinal (within a year) day number.
func ordinalInYear(year int, month time.Month, day int) (dayNo int) {
	return DateToJulian(year, month, day) - DateToJulian(year, 1, 1) + 1
}

// ISOWeekday returns the ISO 8601 weekday number of given day.
// (1 = Mon, 2 = Tue,.. 7 = Sun)
//
// This is different from Go's standard time.Weekday.
func ISOWeekday(year int, month time.Month, day int) (weekday int) {
	// Richards, E. G. (2013) pp. 592, 618

	return DateToJulian(year, month, day)%7 + 1
}

// DateToJulian converts a date to a Julian day number.
func DateToJulian(year int, month time.Month, day int) (dayNo int) {
	// Claus Tøndering's Calendar FAQ
	// http://www.tondering.dk/claus/cal/julperiod.php#formula

	if month < 3 {
		year = year - 1
		month = month + 12
	}
	year = year + 4800

	return day + (153*(int(month)-3)+2)/5 + 365*year +
		year/4 - year/100 + year/400 - 32045
}

// JulianToDate converts a Julian day number to a date.
func JulianToDate(dayNo int) (year int, month time.Month, day int) {
	// Richards, E. G. (2013) pp. 585–624

	e := 4*(dayNo+1401+(4*dayNo+274277)/146097*3/4-38) + 3
	h := e%1461/4*5 + 2

	day = h%153/5 + 1
	month = time.Month((h/153+2)%12 + 1)
	year = e/1461 - 4716 + (14-int(month))/12

	return year, month, day
}

// startOffset returns the offset (in days) from the start of a year to
// Monday of the given week. Offset may be negative.
func startOffset(y, week int) (offset int) {
	// This is optimized version of the following:
	//
	// return week*7 - ISOWeekday(y, 1, 4) - 3
	//
	// Uses Tomohiko Sakamoto's algorithm for calculating the weekday.

	y = y - 1
	return week*7 - (y+y/4-y/100+y/400+3)%7 - 4
}

// StartDate returns the starting date (Monday) of the given ISO 8601 week.
func StartDate(wyear, week int) (year int, month time.Month, day int) {
	return JulianToDate(
		DateToJulian(wyear, 1, 1) + startOffset(wyear, week))
}

// FromYearWeek ...
func FromYearWeek(year, month, day int) (wyear, week int) {
	return fromYearWeek(year, time.Month(month), day)
}

func fromYearWeek(year int, month time.Month, day int) (wyear, week int) {
	week = (ordinalInYear(year, month, day) - ISOWeekday(year, month, day) + 10) / 7
	if week < 1 {
		return fromYearWeek(year-1, 12, 31) // last week of preceding year
	}

	if week == 53 &&
		DateToJulian(StartDate(year+1, 1)) <= DateToJulian(year, month, day) {
		return year + 1, 1 // first week of following year
	}
	return year, week
}

// NextMonthDate 获取当前日期的下个月的第一天和最后一天
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

// PrevMonthDate NextMonthDate 获取当前日期的上个月的第一天和最后一天
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

// MonthDate 获取当前日期的当月的第一天和最后一天
func MonthDate(d string) (string, string) {
	now := cast.ToTimeInDefaultLocation(d, TimeZone())
	year, month, _ := now.Date()
	currentLocation := now.Location()
	firstDate := time.Date(year, month, 1, 0, 0, 0, 0, currentLocation)
	lastDate := firstDate.AddDate(0, 1, -1)
	return Date("Y-m-d", firstDate), Date("Y-m-d", lastDate)
}

// Month 获取当前月份前一年的月份 d Y-m-d
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

// Day 获取起止日期中的天 s,e Y-m-d
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

// Hour Day 获取起止日期中的小时 YmdH
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

// ToWeekDay ...
func ToWeekDay(t interface{}) string {
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
	week := y + (y / 4) + (c / 4) - 2*c + ((26 * (m + 1)) / 10) + day - 1
	if week < 0 {
		week = 7 - (-week)%7
	} else {
		week = week % 7
	}
	which_week := int(week)
	return weekday[which_week]
}
