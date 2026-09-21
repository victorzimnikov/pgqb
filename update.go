package pgqb

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type UpdateBuilder struct {
	Builder

	table string

	where  whereCause
	fields []string

	buildCtx buildContext

	err error
}

func (b *Builder) Update(table string) *UpdateBuilder {
	return &UpdateBuilder{
		table:    table,
		buildCtx: buildContext{},
		Builder:  *b,
	}
}

func (b *UpdateBuilder) Where(field string, arg any) *UpdateBuilder {
	if b.err != nil {
		return b
	}

	b.where.add(field, arg)

	return b
}

func (b *UpdateBuilder) WhereNull(field string) *UpdateBuilder {
	if b.err != nil {
		return b
	}

	b.where.add(field, nil)

	return b
}

func (b *UpdateBuilder) WhereOr(build func(*UpdateBuilder)) *UpdateBuilder {
	if b.err != nil {
		return b
	}

	group := new(UpdateBuilder)
	build(group)

	b.where.addOr(group.where)

	return b
}

func (b *UpdateBuilder) Set(field string, value any) *UpdateBuilder {
	if b.err != nil {
		return b
	}

	newField := quoteIdent(field) + " = " + b.buildCtx.bind(value)

	b.fields = append(b.fields, newField)

	return b
}

func (b *UpdateBuilder) SetIncrement(field string, count int) *UpdateBuilder {
	if b.err != nil {
		return b
	}

	quotedField := quoteIdent(field)

	newField := fmt.Sprintf("%s = %s + %s", quotedField, quotedField, b.buildCtx.bind(count))

	b.fields = append(b.fields, newField)

	return b
}

func (b *UpdateBuilder) SetDecrement(field string, count int) *UpdateBuilder {
	if b.err != nil {
		return b
	}

	quotedField := quoteIdent(field)

	newField := fmt.Sprintf("%s = %s - %s", quotedField, quotedField, b.buildCtx.bind(count))

	b.fields = append(b.fields, newField)

	return b
}

func (b *UpdateBuilder) SetExpr(field string, exp Expression) *UpdateBuilder {
	if b.err != nil {
		return b
	}

	if exp.sql == "" {
		b.err = fmt.Errorf("expression for field %q is empty", field)
		return b
	}

	newField := fmt.Sprintf("%s = %s", quoteIdent(field), exp.sql)

	b.fields = append(b.fields, newField)

	return b
}

func (b *UpdateBuilder) Exec() (pgconn.CommandTag, error) {
	if b.err != nil {
		return pgconn.CommandTag{}, b.err
	}

	if len(b.fields) == 0 {
		return pgconn.CommandTag{}, fmt.Errorf("fields is required")
	}

	query := fmt.Sprintf("UPDATE %s SET %s", quoteIdent(b.table), strings.Join(b.fields, ", "))

	execCtx := buildContext{
		args: append([]any(nil), b.buildCtx.args...),
	}

	whereSQL, err := b.where.build(&execCtx)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	if whereSQL == "" {
		return pgconn.CommandTag{}, fmt.Errorf("update requires WHERE clause")
	}

	query += " " + whereSQL

	return b.db.Exec(b.ctx, query, execCtx.args...)
}
