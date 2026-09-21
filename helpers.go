package pgqb

import "strings"

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func quoteIdentSlice(s []string) string {
	values := make([]string, len(s))

	for i, field := range s {
		values[i] = quoteIdent(field)
	}

	return strings.Join(values, ", ")
}

func joinExpressionSlice(s []Expression) string {
	values := make([]string, len(s))

	for i, field := range s {
		values[i] = field.sql
	}

	return strings.Join(values, ", ")
}

func bindSlice(ctx *buildContext, s []any) string {
	values := make([]string, len(s))

	for i, field := range s {
		values[i] = ctx.bind(field)
	}

	return strings.Join(values, ", ")
}
