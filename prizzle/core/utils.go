package prizzle

import (
	"reflect"
	"strconv"
)

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

func extractMutationsFromValuePairsWithInterceptor(
	slotStartsAt int,
	valuePairs SqlValues,
	interceptor func(
		column string,
		slots string,
		values []interface{},
	),
) {
	for key, value := range valuePairs {
		var column = key.String()

		var sqlValue = getSqlValue(value)

		var values []interface{}

		var slots = getCommaSeparatedSlotList(slotStartsAt, sqlValue.Value, func(value interface{}) {
			values = append(values, value)
		})

		slots += sqlValue.Prefix + slots + sqlValue.Suffix

		interceptor(column, slots, values)
	}
}

func getCommaSeparatedSlotList(startsAt int, values interface{}, adapter func(interface{})) string {
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

// returns the left condition if both conditions is empty
func shouldReturnJustOneConditionIfTheOtherIsEmpty(
	left SqlCondition,
	right SqlCondition,
	resolver func(left SqlCondition, right SqlCondition) SqlCondition,
) SqlCondition {
	if left == "" && right == "" {
		return left
	}

	if left != "" && right != "" {
		return resolver(right, left)
	}

	if left == "" {
		return right
	}

	return left
}

func surroundWithParenthesis(str string) string {
	return "(" + str + ")"
}
