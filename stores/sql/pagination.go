package sql

import "github.com/spf13/cast"

type Pagination struct {
	Offset *int64
	Limit  *int64
}

func (obj *Pagination) GetOffset() int64 {
	return cast.ToInt64(obj.Offset)
}

func (obj *Pagination) GetLimit() int64 {
	return cast.ToInt64(obj.Limit)
}
