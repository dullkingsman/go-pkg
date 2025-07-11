package core

import "encoding/json"

type Nullable[T any] struct {
	value *T
}

func (n *Nullable[T]) IsNull() bool {
	return n.value == nil
}

func (n *Nullable[T]) Get() *T {
	return n.value
}

func (n *Nullable[T]) Set(value *T) {
	n.value = value
}

func (n *Nullable[T]) OrElse(defaultValue T) T {
	if n.IsNull() {
		return defaultValue
	}

	return *n.value
}

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.value = nil
		return nil
	}

	var value T

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.value = &value

	return nil
}

func (n *Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.IsNull() {
		return []byte("null"), nil
	}

	return json.Marshal(*n.value)
}
