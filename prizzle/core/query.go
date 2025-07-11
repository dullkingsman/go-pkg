package prizzle

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Deprecated: DO NOT USE THIS FUNCTION!
func (c *Conn) Prepare(query string) (*sql.Stmt, error) { //TODO find a way to remove this
	panic("This function is decorative. It is here because I could not get around the type system of go. IT SHOULD NOT BE USED!")
}

// Deprecated: DO NOT USE THIS FUNCTION!
func (c *Conn) Exec(query string, args ...any) (sql.Result, error) { //TODO find a way to remove this
	panic("This function is decorative. It is here because I could not get around the type system of go. IT SHOULD NOT BE USED!")
}

// Deprecated: DO NOT USE THIS FUNCTION!
func (c *Conn) Query(query string, args ...any) (*sql.Rows, error) { //TODO find a way to remove this
	panic("This function is decorative. It is here because I could not get around the type system of go. IT SHOULD NOT BE USED!")
}

// Deprecated: DO NOT USE THIS FUNCTION!
func (c *Conn) QueryRow(query string, args ...any) *sql.Row { //TODO find a way to remove this
	panic("This function is decorative. It is here because I could not get around the type system of go. IT SHOULD NOT BE USED!")
}

func (s *Stmt) Exec() (sql.Result, error) {
	return s.Stmt.Exec(s.Args...)
}

func (s *Stmt) ExecContext(ctx context.Context) (sql.Result, error) {
	return s.Stmt.ExecContext(ctx, s.Args...)
}

func (s *Stmt) Query() (*sql.Rows, error) {
	return s.Stmt.Query(s.Args...)
}

func (s *Stmt) QueryContext(ctx context.Context) (*sql.Rows, error) {
	return s.Stmt.QueryContext(ctx, s.Args...)
}

func (s *Stmt) QueryRow() *sql.Row {
	return s.Stmt.QueryRow(s.Args...)
}

func (s *Stmt) QueryRowContext(ctx context.Context) *sql.Row {
	return s.Stmt.QueryRowContext(ctx, s.Args...)
}

func (s *Stmt) SetArgs(args ...any) *Stmt {
	s.Args = args
	return s
}

func (c *Conn) NewQuery() *SqlQuery {
	var query = Query()

	query.Client = c

	return query
}

func (client *DB) NewQuery() *SqlQuery {
	var query = Query()

	query.Client = client

	return query
}

func (transactor *Tx) NewQuery() *SqlQuery {
	var query = Query()

	query.Client = transactor

	return query
}

func Query() *SqlQuery {
	return &SqlQuery{
		Args: &[]interface{}{},
	}
}

func (q *SqlQuery) SubQuery() *SqlQuery {
	return &SqlQuery{
		Args: q.Args,
	}
}

func (q *SqlQuery) QueryString() SqlQueryString {
	return SqlQueryString(q.WithStr + " " + q.BaseStr + " " + q.FromStr + " " + q.WhereStr + " " + q.GroupByStr + " " + q.OrderStr + " " + q.LimitStr + " " + q.OffsetStr + " " + q.ReturningStr)
}

func (q *SqlQuery) Build() PreparedSqlQuery {
	return PreparedSqlQuery{
		Client:      q.Client,
		QueryString: q.QueryString().String(),
		Args:        *q.Args,
	}
}

func (q *PreparedSqlQuery) Exec() (sql.Result, error) {
	return q.Client.Exec(q.QueryString, q.Args...)
}

func (q *PreparedSqlQuery) ExecContext(ctx context.Context) (sql.Result, error) {
	return q.Client.ExecContext(ctx, q.QueryString, q.Args...)
}

func (q *PreparedSqlQuery) Query() (*sql.Rows, error) {
	return q.Client.Query(q.QueryString, q.Args...)
}

func (q *PreparedSqlQuery) QueryContext(ctx context.Context) (*sql.Rows, error) {
	return q.Client.QueryContext(ctx, q.QueryString, q.Args...)
}

func (q *PreparedSqlQuery) QueryRow() *sql.Row {
	return q.Client.QueryRow(q.QueryString, q.Args...)
}

func (q *PreparedSqlQuery) QueryRowContext(ctx context.Context) *sql.Row {
	return q.Client.QueryRowContext(ctx, q.QueryString, q.Args...)
}

func (q *PreparedSqlQuery) Prepare() (*Stmt, error) {
	var stmt, err = q.Client.Prepare(q.QueryString)

	if err != nil {
		return nil, err
	}

	return &Stmt{Stmt: stmt, Args: q.Args}, nil
}

func (q *PreparedSqlQuery) PrepareContext(ctx context.Context) (*Stmt, error) {
	var stmt, err = q.Client.PrepareContext(ctx, q.QueryString)

	if err != nil {
		return nil, err
	}

	return &Stmt{Stmt: stmt, Args: q.Args}, nil
}

// WITH ----------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) With(alias SqlName, as SqlQueryString) *SqlQuery {
	var prefix = "\nWITH "

	if q.WithStr != "" {
		prefix = ", "
	}

	q.WithStr += getPrefixedList(
		"",
		prefix,
		[]string{alias.String() + " AS (" + as.String() + "\n)"},
	)

	return q
}

// INSERT --------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) InsertInto(target EmbedsSqlTable, valuePairs SqlValues) *SqlQuery {
	var valuePointers []string

	var columns, _ = extractMutationsFromValuePairsWithInterceptor(
		valuePairs,
		func(column string, value interface{}) {
			var sqlValue = getSqlValue(value)

			var valueIndicators = getCommaSeparatedStringList(len(*q.Args)+1, sqlValue.Value, func(value interface{}) {
				*q.Args = append(*q.Args, value)
			})

			valuePointers = append(valuePointers, fmt.Sprintf("%s%s%s", sqlValue.Prefix, valueIndicators, sqlValue.Suffix))
		},
	)

	if q.BaseStr != "" {
		q.BaseStr = ","
	} else {
		q.BaseStr += getPrefixedList(
			q.BaseStr,
			"\nINSERT INTO"+" "+target.GetSqlTable().String()+" (",
			columns,
		) + ")\nVALUES"
	}

	q.BaseStr += getPrefixedList(
		"",
		"\n\t(",
		valuePointers,
	) + ")"

	return q
}

func (q *SqlQuery) BulkInsertInto(target EmbedsSqlTable, valuePairsList []SqlValues) *SqlQuery {
	if len(valuePairsList) == 0 {
		return q
	}

	var valuePointers []string

	var columns, _ = extractMutationsFromValuePairsWithInterceptor(
		valuePairsList[0],
		func(column string, value interface{}) {
			var sqlValue = getSqlValue(value)

			var valueIndicators = getCommaSeparatedStringList(len(*q.Args)+1, sqlValue.Value, func(value interface{}) {
				*q.Args = append(*q.Args, value)
			})

			valuePointers = append(valuePointers, fmt.Sprintf("%s%s%s", sqlValue.Prefix, valueIndicators, sqlValue.Suffix))
		},
	)

	if q.BaseStr != "" {
		q.BaseStr = ","
	} else {
		q.BaseStr += getPrefixedList(
			q.BaseStr,
			"\nINSERT INTO"+" "+target.GetSqlTable().String()+" (",
			columns,
		) + ")\nVALUES"
	}

	for i, valuePairs := range valuePairsList {
		if i == 0 {
			continue
		}

		valuePointers = []string{}

		extractMutationsFromValuePairsWithInterceptor(
			valuePairs,
			func(column string, value interface{}) {
				var sqlValue = getSqlValue(value)

				var valueIndicators = getCommaSeparatedStringList(len(*q.Args)+1, sqlValue.Value, func(value interface{}) {
					*q.Args = append(*q.Args, value)
				})

				valuePointers = append(valuePointers, fmt.Sprintf("%s%s%s", sqlValue.Prefix, valueIndicators, sqlValue.Suffix))
			},
		)

		q.BaseStr += getPrefixedList(
			"",
			"\n\t(",
			valuePointers,
		) + ")"

		if i != len(valuePairsList)-1 {
			q.BaseStr += ","
		}
	}

	return q
}

func (q *SqlQuery) InsertIntoFromSelect(target EmbedsSqlTable, columns ...SqlName) *SqlQuery {
	q.BaseStr += getPrefixedList(
		q.BaseStr,
		"\nINSERT INTO"+" "+target.GetSqlTable().String()+" (",
		sqlNameListToStringList(columns),
	) + ")"

	return q
}

// UPDATE --------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) Update(target EmbedsSqlTable, valuePairs SqlValues) *SqlQuery {
	var updates []string

	extractMutationsFromValuePairsWithInterceptor(
		valuePairs,
		func(column string, value interface{}) {
			var sqlValue = getSqlValue(value)

			var valueIndicators = getCommaSeparatedStringList(len(*q.Args)+1, sqlValue.Value, func(value interface{}) {
				*q.Args = append(*q.Args, value)
			})

			updates = append(updates, column+"="+fmt.Sprintf("%s%s%s", sqlValue.Prefix, valueIndicators, sqlValue.Suffix))
		},
	)

	q.BaseStr += getPrefixedList(
		q.BaseStr,
		"\nUPDATE"+" "+target.GetSqlTable().String()+" SET ",
		updates,
	)

	return q
}

func (q *SqlQuery) SetFromSubQuery(column SqlName, subQuery SqlQueryString) *SqlQuery {
	q.BaseStr += getPrefixedList(
		q.BaseStr,
		"",
		[]string{column.String() + " = (" + subQuery.String() + ")"},
	)

	return q
}

func (q *SqlQuery) SetFromOperation(column SqlName, operation SqlOperation) *SqlQuery {
	q.BaseStr += getPrefixedList(
		q.BaseStr,
		"",
		[]string{column.String() + " = " + operation.String()},
	)

	return q
}

// UPSERT --------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) OnConflict(columns ...SqlName) *SqlQuery {
	var prefix = "\nON CONFLICT"

	if len(columns) > 0 {
		prefix = prefix + " ("
	}

	q.BaseStr += getPrefixedList(
		"",
		prefix,
		sqlNameListToStringList(columns),
	) + ")"

	return q
}

func (q *SqlQuery) DoUpdate(valuePairs SqlValues) *SqlQuery {
	var updates []string

	extractMutationsFromValuePairsWithInterceptor(
		valuePairs,
		func(column string, value interface{}) {
			var sqlValue = getSqlValue(value)

			var valueIndicators = getCommaSeparatedStringList(len(*q.Args)+1, sqlValue.Value, func(value interface{}) {
				*q.Args = append(*q.Args, value)
			})

			updates = append(updates, column+" = "+fmt.Sprintf("%s%s%s", sqlValue.Prefix, valueIndicators, sqlValue.Suffix))
		},
	)

	q.BaseStr += getPrefixedList(
		"",
		"\nDO UPDATE SET ",
		updates,
	)

	return q
}

func (q *SqlQuery) DoNothing() *SqlQuery {
	q.BaseStr += "\nDO NOTHING "

	return q
}

// DELETE --------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) DeleteFrom(target EmbedsSqlTable) *SqlQuery {
	q.BaseStr += getPrefixedList(
		q.BaseStr,
		"\nDELETE FROM",
		[]string{target.GetSqlTable().String()},
	)

	return q
}

// RETURN --------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) Returning(columns ...SqlName) *SqlQuery {
	q.ReturningStr += getPrefixedList(
		q.ReturningStr,
		"RETURNING ",
		sqlNameListToStringList(columns),
	)

	return q
}

// SELECT --------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) selectNamespacedWithPrefix(prefix string, namespace string, selectables []SqlName) *SqlQuery {
	var _selectables = sqlNameListToStringList(selectables)

	var _prefix = "SELECT"

	if prefix != "" {
		_prefix = prefix
	}

	q.BaseStr += getPrefixedListWithExtractor(
		q.BaseStr,
		_prefix,
		_selectables,
		func(index int, item string) string {
			var tmp = ""

			if index == 0 && len(_selectables) > 1 {
				tmp += "\n\t"
			} else if index == 0 {
				tmp += " "
			}

			if namespace != "" {
				tmp += surroundWithDoubleQuotes(namespace) + "."
			}

			tmp += item

			if index != len(_selectables)-1 {
				tmp += ",\n\t"
			}

			return tmp
		},
	)

	return q
}

func (q *SqlQuery) Select(selectables ...SqlName) *SqlQuery {
	return q.selectNamespacedWithPrefix("\nSELECT", "", selectables)
}

func (q *SqlQuery) SelectNamespaced(namespace string, selectables ...SqlName) *SqlQuery {
	return q.selectNamespacedWithPrefix("\nSELECT", namespace, selectables)
}

func (q *SqlQuery) SelectDistinct(selectables ...SqlName) *SqlQuery {
	return q.selectNamespacedWithPrefix("\nSELECT DISTINCT", "", selectables)
}

func (q *SqlQuery) SelectDistinctNamespaced(namespace string, selectables ...SqlName) *SqlQuery {
	return q.selectNamespacedWithPrefix("\nSELECT DISTINCT", namespace, selectables)
}

func (q *SqlQuery) SelectDistinctOn(differentiator SqlName, selectables ...SqlName) *SqlQuery {
	return q.selectNamespacedWithPrefix("\nSELECT DISTINCT ON ("+differentiator.String()+")", "", selectables)
}

func (q *SqlQuery) SelectDistinctOnNamespaced(namespace string, differentiator SqlName, selectables ...SqlName) *SqlQuery {
	return q.selectNamespacedWithPrefix("\nSELECT DISTINCT ON ("+differentiator.String()+")", namespace, selectables)
}

// AGGREGATE -----------------------------------------------------------------------------------------------------------

func (q *SqlQuery) Count(column ...SqlName) SqlName {
	var tmp = ""

	if len(column) > 0 {
		tmp += "COUNT(" + column[0].String() + ")"
	} else {
		tmp += "COUNT(*)"
	}

	return SqlName(tmp)
}

func (q *SqlQuery) CountDistinct(column ...SqlName) SqlName {
	var tmp = ""

	if len(column) > 0 {
		tmp += "COUNT(DISTINCT " + column[0].String() + ")"
	} else {
		tmp += "COUNT(*)"
	}

	return SqlName(tmp)
}

func (q *SqlQuery) Avg(column SqlName) SqlName {
	return SqlName("AVG(" + column.String() + ")")
}

func (q *SqlQuery) AvgDistinct(column SqlName) SqlName {
	return SqlName("AVG(DISTINCT " + column.String() + ")")
}

func (q *SqlQuery) Sum(column SqlName) SqlName {
	return SqlName("SUM(" + column.String() + ")")
}

func (q *SqlQuery) SumDistinct(column SqlName) SqlName {
	return SqlName("SUM(DISTINCT " + column.String() + ")")
}

func (q *SqlQuery) Min(column SqlName) SqlName {
	return SqlName("MIN(" + column.String() + ")")
}

func (q *SqlQuery) Max(column SqlName) SqlName {
	return SqlName("MAX(" + column.String() + ")")
}

func (q *SqlQuery) Coalesce(items ...interface{}) SqlName {
	var tmp = "COALESCE("

	var numberOfItems = len(items)

	for index, item := range items {
		if value, ok := item.(SqlName); ok {
			tmp += value.String()
		} else if value, ok := item.(EmbedsSqlTable); ok {
			tmp += value.GetSqlTable().String()
		} else if value, ok := item.(SqlTable); ok {
			tmp += value.String()
		} else if value, ok := item.(string); ok {
			tmp += value
		}

		if index < numberOfItems-1 && numberOfItems > 0 {
			tmp += ", "
		}
	}

	return SqlName(tmp + ")")
}

func (q *SqlQuery) FilterWhere(condition SqlCondition) SqlName {
	return SqlName("FILTER(WHERE " + condition.String() + ")")
}

func (q *SqlQuery) ArrayAgg(array SqlName) SqlName {
	return SqlName("ARRAY_AGG(" + array.String() + ")")
}

func (q *SqlQuery) ArrayAggDistinct(array SqlName) SqlName {
	return SqlName("ARRAY_AGG(DISTINCT " + array.String() + ")")
}

func (q *SqlQuery) JsonbAgg(array SqlName) SqlName {
	return SqlName("JSONB_AGG(" + array.String() + ")")
}

func (q *SqlQuery) JsonbAggDistinct(array SqlName) SqlName {
	return SqlName("JSONB_AGG(DISTINCT " + array.String() + ")")
}

// FROM ----------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) From(source EmbedsSqlTable) *SqlQuery {
	if q.FromStr == "" {
		q.FromStr = "\nFROM " + source.GetSqlTable().String()
	} else {
		q.FromStr += ", " + source.GetSqlTable().String()
	}

	return q
}

func (q *SqlQuery) FromC(source SqlName) *SqlQuery {
	if q.FromStr == "" {
		q.FromStr = "\nFROM " + source.String()
	} else {
		q.FromStr += ", " + source.String()
	}

	return q
}

func (q *SqlQuery) FromSubquery(subQuery SqlQueryString, alias ...SqlName) *SqlQuery {
	if q.FromStr == "" {
		q.FromStr = "\nFROM (" + subQuery.String() + ")"

		if len(alias) > 0 {
			q.FromStr += " " + alias[0].String()
		}
	}

	return q
}

// SUB QUERY -----------------------------------------------------------------------------------------------------------

func (q *SqlQuery) Any(array string) SqlName {
	return SqlName("ANY(" + array + ")")
}

func (q *SqlQuery) Array(subquery SqlQueryString) SqlName {
	return SqlName("ARRAY(" + subquery.String() + ")")
}

func (q *SqlQuery) Unnest(column SqlName, as ...SqlName) SqlName {
	var tmp = "UNNEST(" + column.String() + ")"

	if len(as) > 0 {
		if as[0] != "" {
			tmp += "::" + as[0].String()
		}
	}

	return SqlName(tmp)
}

func (q *SqlQuery) UnnestSub(subQuery SqlQueryString, as ...SqlName) SqlName {
	var tmp = "UNNEST(" + subQuery.String() + ")"

	if len(as) > 0 {
		if as[0] != "" {
			tmp += "::" + as[0].String()
		}
	}

	return SqlName(tmp)
}

func (q *SqlQuery) UnnestArray(array []interface{}, as ...SqlName) SqlName {
	var tmp = "UNNEST(ARRAY["

	for i, value := range array {
		*q.Args = append(*q.Args, value)

		tmp += fmt.Sprintf("$%d", len(*q.Args))

		if i != len(array)-1 {
			tmp += ", "
		}
	}

	tmp = tmp + "])"

	if len(as) > 0 {
		if as[0] != "" {
			tmp += "::" + as[0].String()
		}
	}

	return SqlName(tmp)
}

// JOIN ----------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) Join(source EmbedsSqlTable, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nJOIN " + source.GetSqlTable().String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) JoinC(source SqlName, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nJOIN " + source.String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) InnerJoin(source EmbedsSqlTable, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nINNER JOIN " + source.GetSqlTable().String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) InnerJoinC(source SqlName, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nINNER JOIN " + source.String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) LeftJoin(source EmbedsSqlTable, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nLEFT JOIN " + source.GetSqlTable().String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) LeftJoinC(source SqlName, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nLEFT JOIN " + source.String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) RightJoin(source EmbedsSqlTable, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nRIGHT JOIN " + source.GetSqlTable().String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) RightJoinC(source SqlName, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nRIGHT JOIN " + source.String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) FullJoin(source EmbedsSqlTable, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nFULL JOIN " + source.GetSqlTable().String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) FullJoinC(source SqlName, condition SqlCondition) *SqlQuery {
	q.FromStr += "\nFULL JOIN " + source.String() + " ON " + condition.String()

	return q
}

func (q *SqlQuery) CrossJoin(source EmbedsSqlTable) *SqlQuery {
	q.FromStr += "\nCROSS JOIN " + source.GetSqlTable().String()

	return q
}

func (q *SqlQuery) CrossJoinC(source SqlName) *SqlQuery {
	q.FromStr += "\nCROSS JOIN " + source.String()

	return q
}

// FILTER --------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) Where(condition SqlCondition) *SqlQuery {
	q.WhereStr = "\nWHERE " + condition.String()
	return q
}

func (q *SqlQuery) Eq(left SqlName, right interface{}) SqlCondition {
	*q.Args = append(*q.Args, right)
	return SqlCondition(left.String() + " = " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) Eqc(left SqlName, right SqlName) SqlCondition {
	return SqlCondition(left.String() + " = " + right.String())
}

func (q *SqlQuery) NotEq(left SqlName, right interface{}) SqlCondition {
	*q.Args = append(*q.Args, right)
	return SqlCondition(left.String() + " <> " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) NotEqc(left SqlName, right SqlName) SqlCondition {
	return SqlCondition(left.String() + " <> " + right.String())
}

func (q *SqlQuery) Gt(left SqlName, right interface{}) SqlCondition {
	*q.Args = append(*q.Args, right)
	return SqlCondition(left.String() + " > " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) Gtc(left SqlName, right SqlName) SqlCondition {
	return SqlCondition(left.String() + " > " + right.String())
}

func (q *SqlQuery) Gte(left SqlName, right interface{}) SqlCondition {
	*q.Args = append(*q.Args, right)
	return SqlCondition(left.String() + " >= " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) GteC(left SqlName, right SqlName) SqlCondition {
	return SqlCondition(left.String() + " >= " + right.String())
}

func (q *SqlQuery) Lt(left SqlName, right interface{}) SqlCondition {
	*q.Args = append(*q.Args, right)
	return SqlCondition(left.String() + " < " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) Ltc(left SqlName, right SqlName) SqlCondition {
	return SqlCondition(left.String() + " < " + right.String())
}

func (q *SqlQuery) Lte(left SqlName, right interface{}) SqlCondition {
	*q.Args = append(*q.Args, right)
	return SqlCondition(left.String() + " <= " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) LteC(left SqlName, right SqlName) SqlCondition {
	return SqlCondition(left.String() + " <= " + right.String())
}

func (q *SqlQuery) IsNull(column SqlName) SqlCondition {
	return SqlCondition(column.String() + " IS NULL")
}

func (q *SqlQuery) IsNotNull(column SqlName) SqlCondition {
	return SqlCondition(column.String() + " IS NOT NULL")
}

func (q *SqlQuery) InArray(column SqlName, values []interface{}) SqlCondition {
	var tmp = column.String() + " IN ("

	for i, value := range values {
		*q.Args = append(*q.Args, value)
		tmp += fmt.Sprintf("$%d", len(*q.Args))

		if i != len(values)-1 {
			tmp += ", "
		}
	}

	tmp += ")"

	return SqlCondition(tmp)
}

func (q *SqlQuery) InSubArray(column SqlName, subQuery string) SqlCondition {
	return SqlCondition(column.String() + " IN (" + subQuery + ")")
}

func (q *SqlQuery) NotInArray(column SqlName, values []interface{}) SqlCondition {
	var tmp = column.String() + " NOT IN ("

	for i, value := range values {
		*q.Args = append(*q.Args, value)
		tmp += fmt.Sprintf("$%d", len(*q.Args))

		if i != len(values)-1 {
			tmp += ", "
		}
	}

	tmp += ")"

	return SqlCondition(tmp)
}

func (q *SqlQuery) NotInSubArray(column SqlName, subQuery SqlQueryString) SqlCondition {
	return SqlCondition(column.String() + " NOT IN (" + subQuery.String() + ")")
}

func (q *SqlQuery) Exists(subQuery SqlQueryString) SqlCondition {
	return SqlCondition("EXISTS (" + subQuery.String() + ")")
}

func (q *SqlQuery) NotExists(subQuery SqlQueryString) SqlCondition {
	return SqlCondition("NOT EXISTS (" + subQuery.String() + ")")
}

func (q *SqlQuery) Between(column SqlName, left interface{}, right interface{}) SqlCondition {
	*q.Args = append(*q.Args, left)
	*q.Args = append(*q.Args, right)
	return SqlCondition(column.String() + " BETWEEN " + fmt.Sprintf("%d", len(*q.Args)-1) + " AND " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) NotBetween(column SqlName, left interface{}, right interface{}) SqlCondition {
	*q.Args = append(*q.Args, left)
	*q.Args = append(*q.Args, right)
	return SqlCondition(column.String() + " NOT BETWEEN " + fmt.Sprintf("%d", len(*q.Args)-1) + " AND " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) Like(column SqlName, right string) SqlCondition {
	*q.Args = append(*q.Args, right)
	return SqlCondition(column.String() + " LIKE " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) ILike(column SqlName, right string) SqlCondition {
	*q.Args = append(*q.Args, right)
	return SqlCondition(column.String() + " ILIKE " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) NotILike(column SqlName, right string) SqlCondition {
	*q.Args = append(*q.Args, right)
	return SqlCondition(column.String() + " NOT ILIKE " + fmt.Sprintf("$%d", len(*q.Args)))
}

func (q *SqlQuery) ArrayContains(column SqlName, values ...interface{}) SqlCondition {
	var argsSize = len(*q.Args) + 1

	*q.Args = append(*q.Args, values...)

	var arrayValues = ""

	for index, _ := range values {
		arrayValues += fmt.Sprintf("$%d", argsSize+index)

		if index != len(values)-1 {
			arrayValues += ", "
		}
	}

	return SqlCondition(column.String() + " @> ARRAY[" + arrayValues + "]")
}

// COMBINE -------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) Not(condition SqlCondition) SqlCondition {
	if condition == "" {
		return condition
	}

	return SqlCondition("NOT (" + condition.String() + ")")
}

func (q *SqlQuery) And(left SqlCondition, right SqlCondition) SqlCondition {
	if shouldReturnJustOneConditionIfTheOtherIsEmpty(left, right) != nil {
		return *shouldReturnJustOneConditionIfTheOtherIsEmpty(left, right)
	}

	return SqlCondition("(" + left.String() + " AND " + right.String() + ")")
}

func (q *SqlQuery) Or(left SqlCondition, right SqlCondition) SqlCondition {
	if shouldReturnJustOneConditionIfTheOtherIsEmpty(left, right) != nil {
		return *shouldReturnJustOneConditionIfTheOtherIsEmpty(left, right)
	}

	return SqlCondition("(" + left.String() + " OR " + right.String() + ")")
}

func shouldReturnJustOneConditionIfTheOtherIsEmpty(left SqlCondition, right SqlCondition) *SqlCondition {
	if left == "" && right == "" {
		return &left
	}

	if left != "" && right != "" {
		return nil
	}

	if left == "" {
		return &right
	}

	return &left
}

// GROUP ---------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) GroupBy(groupBy ...SqlName) *SqlQuery {
	q.GroupByStr = "\nGROUP BY "

	for _, group := range groupBy {
		q.GroupByStr += group.String()

		if group != groupBy[len(groupBy)-1] {
			q.GroupByStr += ", "
		}
	}

	return q
}

func (q *SqlQuery) Having(condition SqlCondition) *SqlQuery {
	q.GroupByStr = "\nHAVING " + condition.String()
	return q
}

// SORT ----------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) OrderBy(orderBy SqlName, direction SqlOrderDirection) *SqlQuery {
	if q.OrderStr != "" {
		q.OrderStr += ", "
	} else {
		q.OrderStr = "\nORDER BY "
	}

	q.OrderStr += orderBy.String() + " " + direction.String()

	return q
}

// PAGINATE ------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) Limit(limit int64) *SqlQuery {
	q.LimitStr = "\nLIMIT " + strconv.FormatInt(limit, 10)
	return q
}

func (q *SqlQuery) Offset(offset int64) *SqlQuery {
	q.OffsetStr = "\nOFFSET " + strconv.FormatInt(offset, 10)
	return q
}

// MISC ----------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) NamespacedAlias(namespace string, alias string) SqlName {
	return SqlName(namespace + "." + alias)
}

func (q *SqlQuery) Namespaced(namespace string, name SqlName) SqlName {
	return SqlName(namespace + "." + name.String())
}

// UTILS ---------------------------------------------------------------------------------------------------------------

func (q *SqlQuery) Func(name string, arguments ...interface{}) SqlName {
	var tmp = name + "("

	for i, argument := range arguments {
		if value, ok := argument.(SqlName); ok {
			tmp += value.String()
		} else if value, ok := argument.(EmbedsSqlTable); ok {
			tmp += value.GetSqlTable().String()
		} else if value, ok := argument.(SqlTable); ok {
			tmp += value.String()
		} else {
			*q.Args = append(*q.Args, argument)
			tmp += fmt.Sprintf("$%d", len(*q.Args))
		}

		if i != len(arguments)-1 {
			tmp += ", "
		}
	}

	return SqlName(tmp + ")")
}

func (q *SqlQuery) JsonbBuildObject(json SqlJson) SqlName {
	var tmp = "JSONB_BUILD_OBJECT("

	var total = len(json)
	var count = 0

	for key, item := range json {
		tmp += "\n\t\t"

		tmp += "'" + key.String() + "', " + item.String()

		if total > 1 && count < total-1 {
			tmp += ","
		}

		count++
	}

	return SqlName(tmp + "\n\t)")
}

func ToJsonB(column SqlName) SqlName {
	return SqlName("TO_JSONB(" + column.String() + ")")
}

func ArrayToJson(column SqlName) SqlName {
	return SqlName("ARRAY_TO_JSON(" + column.String() + ")")
}

func JsonbArrayElements(column SqlName) SqlName {
	return SqlName("JSONB_ARRAY_ELEMENTS(" + column.String() + ")")
}

func (q *SqlQuery) JsonArray(value interface{}) SqlValue {
	var returnable = SqlValue{
		Prefix: "ARRAY[",
		Value:  []string{},
		Suffix: "]::json[]",
	}

	var _value = reflect.ValueOf(value)

	if _value.Kind() == reflect.String {
		returnable.Value = append(returnable.Value.([]string), _value.String())
		return returnable

	}

	if _value.Kind() == reflect.Ptr && _value.Elem().Kind() == reflect.String {
		returnable.Value = append(returnable.Value.([]string), _value.Elem().String())
		return returnable
	}

	if _value.Kind() != reflect.Slice {
		return returnable
	}

	for index := 0; index < _value.Len(); index++ {
		var item = _value.Index(index).Interface()

		// todo handle error
		var jsonValue, _ = json.Marshal(item)

		returnable.Value = append(returnable.Value.([]string), string(jsonValue))
	}

	return returnable
}

func (q *SqlQuery) CastToJsonB(column SqlName) SqlName {
	return SqlName(column.String() + "::jsonb")
}

func (q *SqlQuery) CastToJsonBArray(column SqlName) SqlName {
	return SqlName(column.String() + "::jsonb[]")
}

func (q *SqlQuery) Lower(value string) SqlName {
	*q.Args = append(*q.Args, value)
	return SqlName("LOWER($" + strconv.Itoa(len(*q.Args)) + ")")
}

func (q *SqlQuery) LowerC(name SqlName) SqlName {
	return SqlName("LOWER(" + name.String() + ")")
}

func (q *SqlQuery) Upper(value string) SqlName {
	*q.Args = append(*q.Args, value)
	return SqlName("UPPER($" + strconv.Itoa(len(*q.Args)) + ")")
}

func (q *SqlQuery) UpperC(name SqlName) SqlName {
	return SqlName("UPPER(" + name.String() + ")")
}

func getCommaSeparatedStringList(startsAt int, values interface{}, adapter func(interface{})) string {
	var returnable = ""

	var _values = reflect.ValueOf(values)

	if _values.Kind() != reflect.Slice {
		adapter(values)
		return "$" + strconv.Itoa(startsAt)
	}

	for index := 0; index < _values.Len(); index++ {
		var value = _values.Index(index).Interface()

		adapter(value)
		returnable += "$" + strconv.Itoa(startsAt+index)

		if index != _values.Len()-1 {
			returnable += ", "
		}
	}

	return returnable
}
