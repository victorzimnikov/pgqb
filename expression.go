package pgqb

type Expression struct {
	sql string
}

func ClockTimestamp() Expression {
	return Expression{
		sql: "clock_timestamp()",
	}
}

func Column(name string) Expression {
	return Expression{
		sql: quoteIdent(name),
	}
}

func ColumnNull(name string) Expression {
	return Expression{
		sql: quoteIdent(name) + " IS NULL",
	}
}

func ColumnNotNull(name string) Expression {
	return Expression{
		sql: quoteIdent(name) + " IS NOT NULL",
	}
}

func CastField(field string, fieldTypes ...FieldType) (Expression, error) {
	result := quoteIdent(field)

	for _, fieldType := range fieldTypes {
		typeSql, err := fieldType.sql()
		if err != nil {
			return Expression{}, err
		}

		result += "::" + typeSql
	}

	return Expression{sql: result}, nil
}
