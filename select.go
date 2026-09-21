package pgqb

import (
	"fmt"
)

type SelectBuilder struct {
	Builder

	table   string
	limit   int
	offset  int
	where   whereCause
	fields  []Expression
	orderBy []string

	lockTables        []string
	lockMode          string
	lockModeBehaviour string

	err error
}

func (b *Builder) Select(table string, fields ...Expression) *SelectBuilder {
	return &SelectBuilder{
		table:   table,
		fields:  fields,
		Builder: *b,
	}
}

func (b *SelectBuilder) Where(field string, arg any) *SelectBuilder {
	if b.err != nil {
		return b
	}

	b.where.add(field, arg)

	return b
}

func (b *SelectBuilder) WhereNull(field string) *SelectBuilder {
	if b.err != nil {
		return b
	}

	b.where.add(field, nil)

	return b
}

func (b *SelectBuilder) WhereOr(build func(*SelectBuilder)) *SelectBuilder {
	if b.err != nil {
		return b
	}

	group := new(SelectBuilder)
	build(group)

	b.where.addOr(group.where)

	return b
}

func (b *SelectBuilder) Limit(limit int) *SelectBuilder {
	if b.err != nil {
		return b
	}

	b.limit = limit

	return b
}

func (b *SelectBuilder) Offset(offset int) *SelectBuilder {
	if b.err != nil {
		return b
	}

	b.offset = offset

	return b
}

func (b *SelectBuilder) OrderBy(field string, direction OrderDirection) *SelectBuilder {
	if b.err != nil {
		return b
	}

	directionSql, err := direction.sql()
	if err != nil {
		b.err = fmt.Errorf(
			"order by %q: %w",
			field,
			err,
		)

		return b
	}

	b.orderBy = append(b.orderBy, quoteIdent(field)+" "+directionSql)

	return b
}

func (b *SelectBuilder) Lock(
	mode LockMode,
) *SelectBuilder {
	if b.err != nil {
		return b
	}

	modeSql, err := mode.sql()
	if err != nil {
		b.err = fmt.Errorf(
			"lock mode %d: %w",
			mode,
			err,
		)

		return b
	}

	b.lockMode = modeSql

	return b
}

func (b *SelectBuilder) LockTables(tables ...string) *SelectBuilder {
	if b.err != nil {
		return b
	}

	if b.lockMode == "" {
		b.err = fmt.Errorf(
			"lock mode not set",
		)

		return b
	}

	b.lockTables = tables

	return b
}

func (b *SelectBuilder) SkipLocked() *SelectBuilder {
	if b.err != nil {
		return b
	}

	if b.lockMode == "" {
		b.err = fmt.Errorf(
			"lock mode not set",
		)

		return b
	}

	if b.lockModeBehaviour != "" {
		b.err = fmt.Errorf("lock behaviour is already set: %s", b.lockModeBehaviour)

		return b
	}

	b.lockModeBehaviour = "SKIP LOCKED"

	return b
}

func (b *SelectBuilder) NoWait() *SelectBuilder {
	if b.err != nil {
		return b
	}

	if b.lockMode == "" {
		b.err = fmt.Errorf(
			"lock mode not set",
		)

		return b
	}

	if b.lockModeBehaviour != "" {
		b.err = fmt.Errorf("lock behaviour is already set: %s", b.lockModeBehaviour)

		return b
	}

	b.lockModeBehaviour = "NOWAIT"

	return b
}

func (b *SelectBuilder) Exec(dest ...any) error {
	if b.err != nil {
		return b.err
	}

	buildCtx := &buildContext{}

	fields := "*"

	if len(b.fields) > 0 {
		fields = joinExpressionSlice(b.fields)
	}

	query := fmt.Sprintf("SELECT %s FROM %s", fields, quoteIdent(b.table))

	whereSQL, err := b.where.build(buildCtx)
	if err != nil {
		return err
	}

	if whereSQL != "" {
		query += " " + whereSQL
	}

	query += b.buildOrderByCause()
	query += b.buildLimitCause(buildCtx)
	query += b.buildOffsetCause(buildCtx)
	query += b.buildLockCause()

	return b.db.QueryRow(b.ctx, query, buildCtx.args...).Scan(dest...)
}
