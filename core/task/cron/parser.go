package cron

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// ParseOption 用于配置解析器包含哪些字段以及启用哪些特性；未包含的字段由解析器填入默认值，
// 选项不改变字段在表达式中的固定顺序。
type ParseOption int

// 解析器字段与特性选项，可按位组合。
const (
	// Second 为秒字段，缺省默认 0。
	Second ParseOption = 1 << iota
	// SecondOptional 为可选秒字段，缺省默认 0。
	SecondOptional
	// Minute 为分钟字段，缺省默认 0。
	Minute
	// Hour 为小时字段，缺省默认 0。
	Hour
	// Dom 为日期字段，缺省默认 *。
	Dom
	// Month 为月份字段，缺省默认 *。
	Month
	// Dow 为星期字段，缺省默认 *。
	Dow
	// DowOptional 为可选星期字段，缺省默认 *。
	DowOptional
	// Descriptor 允许识别 @monthly、@weekly 等预定义描述符。
	Descriptor
)

var places = []ParseOption{
	Second,
	Minute,
	Hour,
	Dom,
	Month,
	Dow,
}

var defaults = []string{
	"0",
	"0",
	"0",
	"*",
	"*",
	"*",
}

// Parser 是可配置的 cron 表达式解析器。
type Parser struct {
	options ParseOption
}

// NewParser 按给定选项创建解析器。
//
// 同时配置多个可选字段（SecondOptional、DowOptional）时无法可靠推断缺省项，会直接 panic。
//
// 示例：
//
//	// 标准五段表达式，不含描述符
//	specParser := NewParser(Minute | Hour | Dom | Month | Dow)
//	sched, err := specParser.Parse("0 0 15 */3 *")
//
//	// 仅日期、月份、星期三段
//	specParser := NewParser(Dom | Month | Dow)
//	sched, err := specParser.Parse("15 */3 *")
//
//	// 星期字段可选
//	specParser := NewParser(Dom | Month | DowOptional)
//	sched, err := specParser.Parse("15 */3")
func NewParser(options ParseOption) Parser {
	optionals := 0
	if options&DowOptional > 0 {
		optionals++
	}
	if options&SecondOptional > 0 {
		optionals++
	}
	if optionals > 1 {
		panic("multiple optionals may not be configured")
	}
	return Parser{options}
}

// Parse 将调度表字符串解析为 Schedule；表达式非法时返回描述性错误。
//
// 接受由 NewParser 配置的 crontab 表达式格式，也接受 @hourly 等描述符；
// 表达式前缀可使用 CRON_TZ=时区 或 TZ=时区 指定该调度表的时区。
func (p Parser) Parse(spec string) (Schedule, error) {
	if len(spec) == 0 {
		return nil, fmt.Errorf("empty spec string")
	}

	// 解析前缀中指定的时区。
	var loc = time.Local
	if strings.HasPrefix(spec, "TZ=") || strings.HasPrefix(spec, "CRON_TZ=") {
		var err error
		i := strings.Index(spec, " ")
		eq := strings.Index(spec, "=")
		if loc, err = time.LoadLocation(spec[eq+1 : i]); err != nil {
			return nil, fmt.Errorf("provided bad location %s: %v", spec[eq+1:i], err)
		}
		spec = strings.TrimSpace(spec[i:])
	}

	// 已启用描述符选项时，先按预定义描述符处理。
	if strings.HasPrefix(spec, "@") {
		if p.options&Descriptor == 0 {
			return nil, fmt.Errorf("parser does not accept descriptors: %v", spec)
		}
		return parseDescriptor(spec, loc)
	}

	fields := strings.Fields(spec)

	// 校验字段数量并为缺省/可选字段补齐默认值。
	var err error
	fields, err = normalizeFields(fields, p.options)
	if err != nil {
		return nil, err
	}

	field := func(field string, r bounds) uint64 {
		if err != nil {
			return 0
		}
		var bits uint64
		bits, err = getField(field, r)
		return bits
	}

	var (
		second     = field(fields[0], seconds)
		minute     = field(fields[1], minutes)
		hour       = field(fields[2], hours)
		dayofmonth = field(fields[3], dom)
		month      = field(fields[4], months)
		dayofweek  = field(fields[5], dow)
	)
	if err != nil {
		return nil, err
	}

	return &SpecSchedule{
		Second:   second,
		Minute:   minute,
		Hour:     hour,
		Dom:      dayofmonth,
		Month:    month,
		Dow:      dayofweek,
		Location: loc,
	}, nil
}

// normalizeFields 将提交的字段补齐为完整六段，为未启用字段填入默认值，
// 同时校验字段数量与配置选项是否匹配。
func normalizeFields(fields []string, options ParseOption) ([]string, error) {
	// 可选字段按启用处理，并统计可选字段数量。
	optionals := 0
	if options&SecondOptional > 0 {
		options |= Second
		optionals++
	}
	if options&DowOptional > 0 {
		options |= Dow
		optionals++
	}
	if optionals > 1 {
		return nil, fmt.Errorf("multiple optionals may not be configured")
	}

	max := 0
	for _, place := range places {
		if options&place > 0 {
			max++
		}
	}
	min := max - optionals

	if count := len(fields); count < min || count > max {
		if min == max {
			return nil, fmt.Errorf("expected exactly %d fields, found %d: %s", min, count, fields)
		}
		return nil, fmt.Errorf("expected %d to %d fields, found %d: %s", min, max, count, fields)
	}

	// 未提供可选字段时补入对应默认值。
	if min < max && len(fields) == min {
		switch {
		case options&DowOptional > 0:
			// TODO(owner): 优化默认值下标的硬编码访问方式。
			fields = append(fields, defaults[5])
		case options&SecondOptional > 0:
			fields = append([]string{defaults[0]}, fields...)
		default:
			return nil, fmt.Errorf("unknown optional field")
		}
	}

	// 为未启用的字段填入默认值。
	n := 0
	expandedFields := make([]string, len(places))
	copy(expandedFields, defaults)
	for i, place := range places {
		if options&place > 0 {
			expandedFields[i] = fields[n]
			n++
		}
	}
	return expandedFields, nil
}

var standardParser = NewParser(
	Minute | Hour | Dom | Month | Dow | Descriptor,
)

// ParseStandard 按标准五段格式（分、时、日、月、星期）解析调度表，顺序不可调换；
// 表达式非法时返回描述性错误。
//
// 同时接受标准 crontab 表达式（如 "* * * * ?"）与预定义描述符（如 "@midnight"、"@every 1h30m"）。
func ParseStandard(standardSpec string) (Schedule, error) {
	return standardParser.Parse(standardSpec)
}

// getField 解析逗号分隔的取值列表，返回各命中值对应的位集合；任一分段非法时返回错误。
func getField(field string, r bounds) (uint64, error) {
	var bits uint64
	ranges := strings.FieldsFunc(field, func(r rune) bool { return r == ',' })
	for _, expr := range ranges {
		bit, err := getRange(expr, r)
		if err != nil {
			return bits, err
		}
		bits |= bit
	}
	return bits, nil
}

// getRange 解析单个取值表达式，支持 number、number-number 与 number-number/step 形式，
// 返回对应位集合；表达式非法时返回错误。
func getRange(expr string, r bounds) (uint64, error) {
	var (
		start, end, step uint
		rangeAndStep     = strings.Split(expr, "/")
		lowAndHigh       = strings.Split(rangeAndStep[0], "-")
		singleDigit      = len(lowAndHigh) == 1
		err              error
	)

	var extra uint64
	if lowAndHigh[0] == "*" || lowAndHigh[0] == "?" {
		start = r.min
		end = r.max
		extra = starBit
	} else {
		start, err = parseIntOrName(lowAndHigh[0], r.names)
		if err != nil {
			return 0, err
		}
		switch len(lowAndHigh) {
		case 1:
			end = start
		case 2:
			end, err = parseIntOrName(lowAndHigh[1], r.names)
			if err != nil {
				return 0, err
			}
		default:
			return 0, fmt.Errorf("too many hyphens: %s", expr)
		}
	}

	switch len(rangeAndStep) {
	case 1:
		step = 1
	case 2:
		step, err = mustParseInt(rangeAndStep[1])
		if err != nil {
			return 0, err
		}

		// 单独的 N/step 等价于 N-max/step。
		if singleDigit {
			end = r.max
		}
		if step > 1 {
			extra = 0
		}
	default:
		return 0, fmt.Errorf("too many slashes: %s", expr)
	}

	if start < r.min {
		return 0, fmt.Errorf("beginning of range (%d) below minimum (%d): %s", start, r.min, expr)
	}
	if end > r.max {
		return 0, fmt.Errorf("end of range (%d) above maximum (%d): %s", end, r.max, expr)
	}
	if start > end {
		return 0, fmt.Errorf("beginning of range (%d) beyond end of range (%d): %s", start, end, expr)
	}
	if step == 0 {
		return 0, fmt.Errorf("step of range should be a positive number: %s", expr)
	}

	return getBits(start, end, step) | extra, nil
}

// parseIntOrName 解析表达式中的数值；names 中存在别名时按别名映射取值。
func parseIntOrName(expr string, names map[string]uint) (uint, error) {
	if names != nil {
		if namedInt, ok := names[strings.ToLower(expr)]; ok {
			return namedInt, nil
		}
	}
	return mustParseInt(expr)
}

// mustParseInt 将表达式解析为非负整数，解析失败或为负数时返回错误。
func mustParseInt(expr string) (uint, error) {
	num, err := strconv.Atoi(expr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int from %s: %s", expr, err)
	}
	if num < 0 {
		return 0, fmt.Errorf("negative number (%d) not allowed: %s", num, expr)
	}

	return uint(num), nil
}

// getBits 置位区间 [min, max] 内按 step 步进的所有位。
func getBits(min, max, step uint) uint64 {
	var bits uint64

	if step == 1 {
		return ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	}

	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}

// all 置位给定区间内的全部位，并附带星号标记位。
func all(r bounds) uint64 {
	return getBits(r.min, r.max, 1) | starBit
}

// parseDescriptor 将预定义描述符转换为对应调度表；无法识别时返回错误。
func parseDescriptor(descriptor string, loc *time.Location) (Schedule, error) {
	switch descriptor {
	case "@yearly", "@annually":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     1 << hours.min,
			Dom:      1 << dom.min,
			Month:    1 << months.min,
			Dow:      all(dow),
			Location: loc,
		}, nil

	case "@monthly":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     1 << hours.min,
			Dom:      1 << dom.min,
			Month:    all(months),
			Dow:      all(dow),
			Location: loc,
		}, nil

	case "@weekly":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     1 << hours.min,
			Dom:      all(dom),
			Month:    all(months),
			Dow:      1 << dow.min,
			Location: loc,
		}, nil

	case "@daily", "@midnight":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     1 << hours.min,
			Dom:      all(dom),
			Month:    all(months),
			Dow:      all(dow),
			Location: loc,
		}, nil

	case "@hourly":
		return &SpecSchedule{
			Second:   1 << seconds.min,
			Minute:   1 << minutes.min,
			Hour:     all(hours),
			Dom:      all(dom),
			Month:    all(months),
			Dow:      all(dow),
			Location: loc,
		}, nil

	}

	const every = "@every "
	if strings.HasPrefix(descriptor, every) {
		duration, err := time.ParseDuration(descriptor[len(every):])
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration %s: %s", descriptor, err)
		}
		return Every(duration), nil
	}

	return nil, fmt.Errorf("unrecognized descriptor: %s", descriptor)
}
