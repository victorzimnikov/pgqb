package pgqb

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type InsertBuilder struct {
	Builder

	table string

	fields []string

	conflictType   InsertConflictType
	conflictFields []string

	err error
}

func (b *Builder) Insert(table string, fields ...string) *InsertBuilder {
	return &InsertBuilder{
		table:   table,
		fields:  fields,
		Builder: *b,
	}
}

func (b *InsertBuilder) ConflictNothing(fields ...string) *InsertBuilder {
	if b.err != nil {
		return b
	}

	b.conflictType = DoNothingConflict
	b.conflictFields = append(b.conflictFields, fields...)

	return b
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
	query += b.buildConflictClause()

	return b.db.Exec(b.ctx, query, buildCtx.args...)
}

func (b *InsertBuilder) buildConflictClause() string {
	if b.conflictType == 0 {
		return ""
	}

	if len(b.conflictFields) == 0 {
		return " ON CONFLICT DO NOTHING"
	}

	fields := make([]string, len(b.conflictFields))

	for i, field := range b.conflictFields {
		fields[i] = quoteIdent(field)
	}

	return fmt.Sprintf(" ON CONFLICT (%s) DO NOTHING", strings.Join(fields, ", "))
}
