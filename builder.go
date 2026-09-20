package pgqb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type buildContext struct {
	args []any
}

func (c *buildContext) bind(value any) string {
	c.args = append(c.args, value)

	return fmt.Sprintf("$%d", len(c.args))
}

type DBTX interface {
	Exec(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	QueryRow(
		ctx context.Context,
		sql string,
		arguments ...any,
	) pgx.Row
}

type Builder struct {
	db  DBTX
	ctx context.Context
}

func NewBuilder(ctx context.Context, db DBTX) *Builder {
	return &Builder{
		ctx: ctx,
		db:  db,
	}
}
