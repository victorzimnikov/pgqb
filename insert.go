package pgqb

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type InsertBuilder struct {
	Builder

	table string

	fields []string

	err error
}

func (b *Builder) Insert(table string, fields ...string) *InsertBuilder {
	return &InsertBuilder{
		table:   table,
		fields:  fields,
		Builder: *b,
	}
}

func (b *InsertBuilder) Exec(values ...any) (pgconn.CommandTag, error) {
	if b.err != nil {
		return pgconn.CommandTag{}, b.err
	}

	if len(b.fields) != len(values) {
		return pgconn.CommandTag{}, fmt.Errorf("values count not equal fields count")
	}

	if len(b.fields) == 0 {
		return pgconn.CommandTag{}, fmt.Errorf("fields is required")
	}

	fields := ""

	if len(b.fields) > 0 {
		for _, field := range b.fields {
			if fields == "" {
				fields = quoteIdent(field)
			} else {
				fields += ", " + quoteIdent(field)
			}
		}
	}

	buildCtx := &buildContext{}

	valuesString := ""

	for _, value := range values {
		if valuesString == "" {
			valuesString = buildCtx.bind(value)
		} else {
			valuesString += ", " + buildCtx.bind(value)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", quoteIdent(b.table), fields, valuesString)

	return b.db.Exec(b.ctx, query, buildCtx.args...)
}
