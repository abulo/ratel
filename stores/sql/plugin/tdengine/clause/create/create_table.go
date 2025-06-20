package create

import (
	"gorm.io/gorm/clause"
)

type CreateTable struct {
	tables []*Table
}

// NewCreateTable Create table clause
func NewCreateTable(tables ...*Table) *CreateTable {
	return &CreateTable{tables: tables}
}

// AddTables Add tables to clause
func (c *CreateTable) AddTables(tables ...*Table) *CreateTable {
	c.tables = append(c.tables, tables...)
	return c
}

func (CreateTable) Name() string {
	return "CREATE TABLE"
}

func (c CreateTable) Build(builder clause.Builder) {
	for _, tb := range c.tables {
		tb.Build(builder)
	}
}

// MergeClause merge CREATE TABLE by clauses
func (c CreateTable) MergeClause(clause *clause.Clause) {
	clause.Name = ""
	clause.Expression = c
}
