package sql

type Pagination struct {
	Offset *int64
	Limit  *int64
}

func (obj *Pagination) GetOffset() int64 {
	return *obj.Offset
}

func (obj *Pagination) GetLimit() int64 {
	return *obj.Limit
}
