package sql

// Epr ...
type Epr struct {
	value string
}

// NewEpr ...
func NewEpr(value string) Epr {
	return Epr{value: value}
}

// ToString ...
func (e Epr) ToString() string {
	return e.value
}
