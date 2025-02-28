package sql

import (
	"errors"

	"github.com/spf13/cast"
)

type Scope struct {
	Start *string
	End   *string
}

func (obj *Scope) GetStartString() (string, error) {
	if obj.Start == nil {
		return "", errors.New("start is nil")
	}
	return cast.ToString(obj.Start), nil
}

func (obj *Scope) GetEndString() (string, error) {
	if obj.End == nil {
		return "", errors.New("end is nil")
	}
	return cast.ToString(obj.End), nil
}

func (obj *Scope) GetStartInt64() (int64, error) {
	if obj.Start == nil {
		return 0, errors.New("start is nil")
	}
	return cast.ToInt64(obj.Start), nil
}

func (obj *Scope) GetEndInt64() (int64, error) {
	if obj.End == nil {
		return 0, errors.New("end is nil")
	}
	return cast.ToInt64(obj.End), nil
}

func (obj *Scope) GetStartFloat64() (float64, error) {
	if obj.Start == nil {
		return 0, errors.New("start is nil")
	}
	return cast.ToFloat64(obj.Start), nil
}

func (obj *Scope) GetEndFloat64() (float64, error) {
	if obj.End == nil {
		return 0, errors.New("end is nil")
	}
	return cast.ToFloat64(obj.End), nil
}
