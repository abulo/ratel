package query

type Epr struct {
	value string
}

func NewEpr(value string) Epr {
	return Epr{value: value}
}
func (e Epr) ToString() string {
	return e.value
}
