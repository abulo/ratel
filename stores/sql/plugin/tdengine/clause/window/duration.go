package window

import (
	"errors"
	"strconv"
	"time"
)

type UnitType string

// u(微秒)、a(毫秒)、s(秒)、m(分)、h(小时)、d(天)、w(周) n(自然月) 和 y(自然年)
const (
	Microsecond UnitType = "u"
	Millisecond UnitType = "a"
	Second      UnitType = "s"
	Minute      UnitType = "m"
	Hour        UnitType = "h"
	Day         UnitType = "d"
	Week        UnitType = "w"
	Month       UnitType = "n"
	Year        UnitType = "y"
)

var durationMap = map[UnitType]struct{}{
	Microsecond: {},
	Millisecond: {},
	Second:      {},
	Minute:      {},
	Hour:        {},
	Day:         {},
	Week:        {},
	Month:       {},
	Year:        {},
}

type Duration struct {
	Value uint64
	Unit  UnitType
}

func NewDuration(duration time.Duration) (*Duration, error) {
	if duration <= 0 {
		return nil, errors.New("duration does not allow negative numbers")
	}
	return &Duration{
		Value: uint64(duration.Microseconds()),
		Unit:  Microsecond,
	}, nil
}

func ParseDuration(durationString string) (*Duration, error) {
	if len(durationString) < 2 {
		return nil, errors.New("parse duration failure")
	}
	unit := UnitType(durationString[len(durationString)-1:])
	_, valid := durationMap[unit]
	if !valid {
		return nil, errors.New("unit not valid")
	}
	value := durationString[:len(durationString)-1]
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Duration{
		Value: v,
		Unit:  unit,
	}, nil
}
