package pgqb

import (
	"fmt"
	"strings"
)

func (b *SelectBuilder) buildOrderByCause() string {
	if len(b.orderBy) == 0 {
		return ""
	}

	return " ORDER BY " + strings.Join(b.orderBy, ", ")
}

func (b *SelectBuilder) buildLimitCause(ctx *buildContext) string {
	if b.limit == 0 {
		return ""
	}

	return fmt.Sprintf(" LIMIT %s", ctx.bind(b.limit))
}

func (b *SelectBuilder) buildOffsetCause(ctx *buildContext) string {
	if b.offset == 0 {
		return ""
	}

	return fmt.Sprintf(" OFFSET %s", ctx.bind(b.offset))
}

func (b *SelectBuilder) buildLockCause() string {
	if b.lockMode == "" {
		return ""
	}

	query := ""

	if len(b.lockTables) > 0 {
		query += fmt.Sprintf(
			" FOR %s OF %s",
			b.lockMode,
			quoteIdentSlice(b.lockTables),
		)
	} else {
		query += fmt.Sprintf(" FOR %s", b.lockMode)
	}

	if b.lockModeBehaviour != "" {
		query += " " + b.lockModeBehaviour
	}

	return query
}
