package window

import (
	"strconv"

	"gorm.io/gorm/clause"
)

const (
	SESSION = iota + 1
	STATE
	INTERVAL
)

// [SESSION(ts_col, tol_val)]
// [STATE_WINDOW(col)]
// [INTERVAL(interval_val [, interval_offset]) [SLIDING sliding_val]]

type Window struct {
	windowType  int
	tsColumn    string
	stateColumn string
	duration    *Duration
	offset      *Duration
	sliding     *Duration
}

// SetSessionWindow create a session window [SESSION(ts_col, tol_val)]
func SetSessionWindow(tsColumn string, duration Duration) Window {
	return Window{windowType: SESSION, tsColumn: tsColumn, duration: &duration}
}

// SetStateWindow create a state window [STATE_WINDOW(col)]
func SetStateWindow(column string) Window {
	return Window{windowType: STATE, stateColumn: column}
}

// SetInterval create an interval window [INTERVAL(interval_val [, interval_offset]) [SLIDING sliding_val]]
func SetInterval(duration Duration) Window {
	return Window{windowType: INTERVAL, duration: &duration}
}

// SetOffset set offset to interval window
func (sc Window) SetOffset(offset Duration) Window {
	if sc.windowType == INTERVAL {
		sc.offset = &offset
	}
	return sc
}

// SetSliding set sliding to interval window
func (sc Window) SetSliding(sliding Duration) Window {
	if sc.windowType == INTERVAL {
		sc.sliding = &sliding
	}
	return sc
}

func (sc Window) Build(builder clause.Builder) {
	switch sc.windowType {
	case SESSION:
		builder.WriteString("SESSION(")
		builder.WriteQuoted(sc.tsColumn)
		builder.WriteByte(',')
		builder.WriteString(strconv.FormatUint(sc.duration.Value, 10))
		builder.WriteString(string(sc.duration.Unit))
		builder.WriteByte(')')
	case STATE:
		builder.WriteString("STATE_WINDOW(")
		builder.WriteQuoted(sc.stateColumn)
		builder.WriteByte(')')
	case INTERVAL:
		builder.WriteString("INTERVAL(")
		builder.WriteString(strconv.FormatUint(sc.duration.Value, 10))
		builder.WriteString(string(sc.duration.Unit))
		if sc.offset != nil {
			builder.WriteByte(',')
			builder.WriteString(strconv.FormatUint(sc.offset.Value, 10))
			builder.WriteString(string(sc.offset.Unit))
		}
		builder.WriteByte(')')
		if sc.sliding != nil {
			builder.WriteString(" SLIDING(")
			builder.WriteString(strconv.FormatUint(sc.sliding.Value, 10))
			builder.WriteString(string(sc.sliding.Unit))
			builder.WriteByte(')')
		}
	}
}

func (sc Window) Name() string {
	return "WINDOW"
}

func (sc Window) MergeClause(c *clause.Clause) {
	c.Name = ""
	c.Expression = sc
}
