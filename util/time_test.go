package util

import (
	"reflect"
	"testing"
	"time"
	_ "time/tzdata"
)

func TestDateTimeParse(t *testing.T) {
	type args struct {
		st string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DateTimeParse(tt.args.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("DateTimeParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DateTimeParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDuring(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDuring(tt.args.t); got != tt.want {
				t.Errorf("FormatDuring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHourDiffer(t *testing.T) {
	type args struct {
		startTime interface{}
		endTime   interface{}
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHourDiffer(tt.args.startTime, tt.args.endTime); got != tt.want {
				t.Errorf("GetHourDiffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeZone(t *testing.T) {
	tests := []struct {
		name string
		want *time.Location
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeZone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeZone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Time(); got != tt.want {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Duration(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNow(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Now(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Now() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Timestamp(); got != tt.want {
				t.Errorf("Timestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToTime(t *testing.T) {
	type args struct {
		format  string
		strtime string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToTime(tt.args.format, tt.args.strtime)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StrToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate(t *testing.T) {
	type args struct {
		format string
		ts     []time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Date(tt.args.format, tt.args.ts...); got != tt.want {
				t.Errorf("Date() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDate(t *testing.T) {
	type args struct {
		month int
		day   int
		year  int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDate(tt.args.month, tt.args.day, tt.args.year); got != tt.want {
				t.Errorf("CheckDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSleep(t *testing.T) {
	type args struct {
		t int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Sleep(tt.args.t)
		})
	}
}

func TestUSleep(t *testing.T) {
	type args struct {
		t int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			USleep(tt.args.t)
		})
	}
}

func TestUnixTimeFormatDate(t *testing.T) {
	type args struct {
		str interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnixTimeFormatDate(tt.args.str); got != tt.want {
				t.Errorf("UnixTimeFormatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDateTime(t *testing.T) {
	type args struct {
		str interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDateTime(tt.args.str); got != tt.want {
				t.Errorf("FormatDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	type args struct {
		str interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDate(tt.args.str); got != tt.want {
				t.Errorf("FormatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ordinalInYear(t *testing.T) {
	type args struct {
		year  int
		month time.Month
		day   int
	}
	tests := []struct {
		name      string
		args      args
		wantDayNo int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDayNo := ordinalInYear(tt.args.year, tt.args.month, tt.args.day); gotDayNo != tt.wantDayNo {
				t.Errorf("ordinalInYear() = %v, want %v", gotDayNo, tt.wantDayNo)
			}
		})
	}
}

func TestISOWeekday(t *testing.T) {
	type args struct {
		year  int
		month time.Month
		day   int
	}
	tests := []struct {
		name        string
		args        args
		wantWeekday int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWeekday := ISOWeekday(tt.args.year, tt.args.month, tt.args.day); gotWeekday != tt.wantWeekday {
				t.Errorf("ISOWeekday() = %v, want %v", gotWeekday, tt.wantWeekday)
			}
		})
	}
}

func TestDateToJulian(t *testing.T) {
	type args struct {
		year  int
		month time.Month
		day   int
	}
	tests := []struct {
		name      string
		args      args
		wantDayNo int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDayNo := DateToJulian(tt.args.year, tt.args.month, tt.args.day); gotDayNo != tt.wantDayNo {
				t.Errorf("DateToJulian() = %v, want %v", gotDayNo, tt.wantDayNo)
			}
		})
	}
}

func TestJulianToDate(t *testing.T) {
	type args struct {
		dayNo int
	}
	tests := []struct {
		name      string
		args      args
		wantYear  int
		wantMonth time.Month
		wantDay   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotYear, gotMonth, gotDay := JulianToDate(tt.args.dayNo)
			if gotYear != tt.wantYear {
				t.Errorf("JulianToDate() gotYear = %v, want %v", gotYear, tt.wantYear)
			}
			if !reflect.DeepEqual(gotMonth, tt.wantMonth) {
				t.Errorf("JulianToDate() gotMonth = %v, want %v", gotMonth, tt.wantMonth)
			}
			if gotDay != tt.wantDay {
				t.Errorf("JulianToDate() gotDay = %v, want %v", gotDay, tt.wantDay)
			}
		})
	}
}

func Test_startOffset(t *testing.T) {
	type args struct {
		y    int
		week int
	}
	tests := []struct {
		name       string
		args       args
		wantOffset int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOffset := startOffset(tt.args.y, tt.args.week); gotOffset != tt.wantOffset {
				t.Errorf("startOffset() = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}

func TestStartDate(t *testing.T) {
	type args struct {
		wyear int
		week  int
	}
	tests := []struct {
		name      string
		args      args
		wantYear  int
		wantMonth time.Month
		wantDay   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotYear, gotMonth, gotDay := StartDate(tt.args.wyear, tt.args.week)
			if gotYear != tt.wantYear {
				t.Errorf("StartDate() gotYear = %v, want %v", gotYear, tt.wantYear)
			}
			if !reflect.DeepEqual(gotMonth, tt.wantMonth) {
				t.Errorf("StartDate() gotMonth = %v, want %v", gotMonth, tt.wantMonth)
			}
			if gotDay != tt.wantDay {
				t.Errorf("StartDate() gotDay = %v, want %v", gotDay, tt.wantDay)
			}
		})
	}
}

func TestFromYearWeek(t *testing.T) {
	type args struct {
		year  int
		month int
		day   int
	}
	tests := []struct {
		name      string
		args      args
		wantWyear int
		wantWeek  int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWyear, gotWeek := FromYearWeek(tt.args.year, tt.args.month, tt.args.day)
			if gotWyear != tt.wantWyear {
				t.Errorf("FromYearWeek() gotWyear = %v, want %v", gotWyear, tt.wantWyear)
			}
			if gotWeek != tt.wantWeek {
				t.Errorf("FromYearWeek() gotWeek = %v, want %v", gotWeek, tt.wantWeek)
			}
		})
	}
}

func Test_fromYearWeek(t *testing.T) {
	type args struct {
		year  int
		month time.Month
		day   int
	}
	tests := []struct {
		name      string
		args      args
		wantWyear int
		wantWeek  int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWyear, gotWeek := fromYearWeek(tt.args.year, tt.args.month, tt.args.day)
			if gotWyear != tt.wantWyear {
				t.Errorf("fromYearWeek() gotWyear = %v, want %v", gotWyear, tt.wantWyear)
			}
			if gotWeek != tt.wantWeek {
				t.Errorf("fromYearWeek() gotWeek = %v, want %v", gotWeek, tt.wantWeek)
			}
		})
	}
}

func TestNextMonthDate(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NextMonthDate(tt.args.d)
			if got != tt.want {
				t.Errorf("NextMonthDate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NextMonthDate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPrevMonthDate(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := PrevMonthDate(tt.args.d)
			if got != tt.want {
				t.Errorf("PrevMonthDate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PrevMonthDate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMonthDate(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := MonthDate(tt.args.d)
			if got != tt.want {
				t.Errorf("MonthDate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MonthDate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMonth(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Month(tt.args.d)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Month() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Month() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDay(t *testing.T) {
	type args struct {
		s string
		e string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Day(tt.args.s, tt.args.e)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Day() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Day() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHour(t *testing.T) {
	type args struct {
		s string
		e string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Hour(tt.args.s, tt.args.e)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hour() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Hour() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestToWeekDay(t *testing.T) {
	type args struct {
		t interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToWeekDay(tt.args.t); got != tt.want {
				t.Errorf("ToWeekDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
