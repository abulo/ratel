package using

import (
	"gorm.io/gorm/clause"
)

type Using struct {
	stable string
	tags   map[string]any
}

// SetUsing Using clause
func SetUsing(stable string, tags map[string]any) Using {
	return Using{
		stable: stable,
		tags:   tags,
	}
}

func (u Using) Build(builder clause.Builder) {
	_, _ = builder.WriteString("USING ")
	builder.WriteQuoted(u.stable)
	tagNames := make([]string, 0, len(u.tags))
	tagValues := make([]any, 0, len(u.tags))
	for tagName, tagValue := range u.tags {
		tagNames = append(tagNames, tagName)
		tagValues = append(tagValues, tagValue)
	}
	builder.AddVar(builder, tagNames)
	_, _ = builder.WriteString(" TAGS")
	builder.AddVar(builder, tagValues)
}

// AddTag add tag pair to using clause
func (u Using) AddTag(tagName string, tagValue any) Using {
	u.tags[tagName] = tagValue
	return u
}

func (Using) Name() string {
	return "USING"
}

func (u Using) MergeClause(c *clause.Clause) {
	c.Name = ""
	c.Expression = u
}
