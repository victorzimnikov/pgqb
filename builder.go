package pgqb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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
