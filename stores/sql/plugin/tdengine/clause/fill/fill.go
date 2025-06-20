package fill

import (
	"strconv"

	"gorm.io/gorm/clause"
)

type FillType string

const (
	FillNone   FillType = "NONE"
	FillValue  FillType = "VALUE"
	FillPrev   FillType = "PREV"
	FillNull   FillType = "NULL"
	FillLinear FillType = "LINEAR"
	FillNext   FillType = "NEXT"
)

type Fill struct {
	Type  FillType
	Value float64 // only support Type = FillValue.
}

// Build [FILL(fill_mod_and_val)]
func (f Fill) Build(builder clause.Builder) {
	_, _ = builder.WriteString("(")
	_, _ = builder.WriteString(string(f.Type))
	if f.Type == FillValue {
		_ = builder.WriteByte(',')
		_, _ = builder.WriteString(strconv.FormatFloat(f.Value, 'g', -1, 64))
	}
	_ = builder.WriteByte(')')
}

func (f Fill) Name() string {
	return "FILL"
}

func (f Fill) MergeClause(c *clause.Clause) {
	c.Expression = f
}
