package sql

import (
	"errors"

	"github.com/spf13/cast"
)

type Pagination struct {
	Offset *int64
	Limit  *int64
}

func (obj *Pagination) GetOffset() (int64, error) {
	if obj.Offset == nil {
		return 0, errors.New("offset is nil")
	}
	return cast.ToInt64(obj.Offset), nil
}

func (obj *Pagination) GetLimit() (int64, error) {
	if obj.Limit == nil {
		return 0, errors.New("limit is nil")
	}
	return cast.ToInt64(obj.Limit), nil
}
