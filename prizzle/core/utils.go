package prizzle

func getPrefixedListWithExtractor(
	base string,
	prefix string,
	items []string,
	extractor func(index int, item string) string,
) string {
	var _base = base

	if _base != "" {
		_base += ","
	} else {
		_base = prefix
	}

	if len(items) == 0 {
		_base += " *"
	} else {
		for i, column := range items {
			_base += extractor(i, column)
		}
	}

	return _base
}

func getPrefixedList(base string, prefix string, items []string) string {
	return getPrefixedListWithExtractor(
		base,
		prefix,
		items,
		func(index int, item string) string {
			var tmp = item

			if index < len(items)-1 {
				tmp += ", "
			}

			return tmp
		},
	)
}

func sqlNameListToStringList(columns []SqlName) []string {
	var items []string

	for _, column := range columns {
		items = append(items, column.String())
	}

	return items
}

func surroundWithDoubleQuotes(string string) string {
	return "\"" + string + "\""
}

func extractMutationsFromValuePairs(valuePairs SqlValues) ([]string, []interface{}) {
	return extractMutationsFromValuePairsWithInterceptor(valuePairs, func(column string, value interface{}) {})
}

func extractMutationsFromValuePairsWithInterceptor(valuePairs SqlValues, interceptor func(column string, value interface{})) ([]string, []interface{}) {
	var columns []string
	var values []interface{}

	for key, value := range valuePairs {
		var column = key.String()

		columns = append(columns, column)
		values = append(values, value)

		interceptor(column, value)
	}

	return columns, values
}

func getSqlValue(value interface{}) SqlValue {
	switch value.(type) {
	case SqlValue:
		return value.(SqlValue)
	}

	return SqlValue{
		Prefix: "",
		Value:  value,
		Suffix: "",
	}
}
