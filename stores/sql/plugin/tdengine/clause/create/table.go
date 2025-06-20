package create

import (
	"errors"
	"strconv"

	"gorm.io/gorm/clause"
)

type Table struct {
	tableType   TableType
	tableName   string
	ifNotExists bool
	// need by STable OR CTable if Stable is empty
	columns []*Column
	// only need by STable
	tagColumns []*Column
	// only need by CTable if stable is not empty
	stable string
	tags   map[string]any
}

func (t Table) TableType() TableType {
	return t.tableType
}

func (t Table) TableName() string {
	return t.tableName
}

func (t Table) Build(builder clause.Builder) {
	switch t.tableType {
	case CTable:
		_, _ = builder.WriteString("CREATE TABLE ")
	case STable:
		_, _ = builder.WriteString("CREATE STABLE ")
	default:
		_ = builder.AddError(errors.New("Unsupported table type"))
		return
	}
	if t.ifNotExists {
		_, _ = builder.WriteString("IF NOT EXISTS ")
	}
	builder.WriteQuoted(t.tableName)
	if t.tableType == CTable && t.stable != "" {
		_, _ = builder.WriteString(" USING ")
		builder.WriteQuoted(t.stable)

		tagValueList := make([]any, 0, len(t.tags))
		index := 0
		_ = builder.WriteByte('(')
		for tag, tagValue := range t.tags {
			if index > 0 {
				_ = builder.WriteByte(',')
			}
			builder.WriteQuoted(tag)

			tagValueList = append(tagValueList, tagValue)
			index += 1
		}
		_, _ = builder.WriteString(") TAGS ")
		builder.AddVar(builder, tagValueList)
	} else {
		_, _ = builder.WriteString(" (")
		for i, column := range t.columns {
			if i > 0 {
				_ = builder.WriteByte(',')
			}
			column.Build(builder)
		}
		_ = builder.WriteByte(')')
	}
	if t.tableType == STable {
		_, _ = builder.WriteString(" TAGS(")
		for i, tags := range t.tagColumns {
			if i > 0 {
				_ = builder.WriteByte(',')
			}
			tags.Build(builder)
		}
		_ = builder.WriteByte(')')
	}
}

type Column struct {
	Type   ColumnType
	Name   string
	Length uint64
}

func (c *Column) Build(builder clause.Builder) {
	builder.WriteQuoted(c.Name)
	_ = builder.WriteByte(' ')
	_, _ = builder.WriteString(string(c.Type))
	if c.Type == NChar ||
		c.Type == Binary ||
		c.Type == VarChar ||
		c.Type == VarBinary {
		_ = builder.WriteByte('(')
		_, _ = builder.WriteString(strconv.FormatUint(c.Length, 10))
		_ = builder.WriteByte(')')
	}
}

type sTableBuilder struct {
	table *Table
}

func NewSTableBuilder(tableName string) *sTableBuilder {
	return &sTableBuilder{table: &Table{
		tableType:   STable,
		tableName:   tableName,
		ifNotExists: false,
		columns:     []*Column{},
		tagColumns:  []*Column{},
		stable:      "",
		tags:        map[string]any{},
	}}
}

func (b *sTableBuilder) Build() *Table {
	return b.table
}

func (b *sTableBuilder) IfNotExists() *sTableBuilder {
	b.table.ifNotExists = true
	return b
}

func (b *sTableBuilder) Columns(columns ...*Column) *sTableBuilder {
	b.table.columns = append(b.table.columns, columns...)
	return b
}

func (b *sTableBuilder) TagColumns(columns ...*Column) *sTableBuilder {
	b.table.tagColumns = append(b.table.tagColumns, columns...)
	return b
}

type cTableBuilder struct {
	table *Table
}

func NewCTableBuilder(tableName string) *cTableBuilder {
	return &cTableBuilder{
		table: &Table{
			tableType:   CTable,
			tableName:   tableName,
			ifNotExists: false,
			columns:     []*Column{},
			tagColumns:  []*Column{},
			stable:      "",
			tags:        map[string]any{},
		},
	}
}

func (b *cTableBuilder) IfNotExists() *cTableBuilder {
	b.table.ifNotExists = true
	return b
}

func (b *cTableBuilder) Columns(columns ...*Column) *cTableBuilder {
	b.table.columns = append(b.table.columns, columns...)
	return b
}

func (b *cTableBuilder) Build() *Table {
	return b.table
}

func (b *cTableBuilder) BuildWithSTable(stable string, tags map[string]any) *Table {
	b.table.stable = stable
	b.table.tags = tags
	return b.table
}
