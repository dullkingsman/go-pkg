package utils

import "reflect"

func PtrOf[T any](value T) *T {
	return &value
}

func ValueOr[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}

	return *value
}

func SafeCastValue[T any](value interface{}) *T {
	if value == nil {
		return nil
	}

	if v, ok := value.(*interface{}); ok {
		if v == nil {
			return nil
		}

		value = *v
	}

	if v, ok := value.(T); ok {
		return &v
	}

	if v, ok := value.(*T); ok {
		return v
	}

	if v, ok := value.(**T); ok && v != nil {
		return *v
	}

	return nil
}

func UnderlyingValueIsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	var val = reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan, reflect.Func:
		return val.IsNil()
	default:
		return false
	}
}
